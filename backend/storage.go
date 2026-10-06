package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// storageConfig menyimpan konfigurasi object storage S3-compatible
// (Cloudflare R2) untuk membuat presigned URL.
type storageConfig struct {
	bucket        string
	presignClient *s3.PresignClient
	publicBaseURL string
	presignTTL    time.Duration
}

type presignUploadRequest struct {
	ContentType string `json:"content_type"`
	Kind        string `json:"kind"`
}

type presignUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	PublicURL string `json:"public_url"`
	ExpiresIn int    `json:"expires_in"`
}

// newStorageConfig membaca konfigurasi R2 dari environment. Bila kredensial
// tidak lengkap, mengembalikan (nil, nil) agar fitur upload dinonaktifkan
// tanpa menggagalkan startup.
func newStorageConfig(ctx context.Context) (*storageConfig, error) {
	accountID := strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID"))
	accessKey := strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID"))
	secretKey := strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY"))
	bucket := strings.TrimSpace(os.Getenv("R2_BUCKET"))
	publicBaseURL := strings.TrimSpace(
		os.Getenv("R2_PUBLIC_BASE_URL"),
	)

	if publicBaseURL == "" {
		publicBaseURL = strings.TrimSpace(os.Getenv("R2_PUBLIC_URL"))
	}

	if accountID == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, nil
	}

	endpoint := fmt.Sprintf(
		"https://%s.r2.cloudflarestorage.com",
		accountID,
	)

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("memuat konfigurasi storage: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &storageConfig{
		bucket:        bucket,
		presignClient: s3.NewPresignClient(client),
		publicBaseURL: strings.TrimSuffix(publicBaseURL, "/"),
		presignTTL:    15 * time.Minute,
	}, nil
}

// allowedUploadContentTypes memetakan MIME type yang diizinkan ke ekstensi.
var allowedUploadContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/heic": ".heic",
	"image/heif": ".heif",
}

// objectPublicURL membentuk URL publik objek berdasarkan object key. Kosong
// bila public base URL tidak dikonfigurasi atau key kosong.
func (s *storageConfig) objectPublicURL(objectKey string) string {
	if s == nil || s.publicBaseURL == "" || objectKey == "" {
		return ""
	}

	return s.publicBaseURL + "/" + objectKey
}

func randomHex(bytes int) (string, error) {
	buffer := make([]byte, bytes)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}

func (app *application) presignUploadHandler(w http.ResponseWriter, r *http.Request) {
	if app.storage == nil {
		http.Error(
			w,
			"object storage belum dikonfigurasi",
			http.StatusServiceUnavailable,
		)
		return
	}

	user, err := app.getAuthenticatedUser(r)
	if err == errUnauthenticated {
		http.Error(w, "belum login", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi user upload: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	var input presignUploadRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	contentType := strings.ToLower(strings.TrimSpace(input.ContentType))

	extension, allowed := allowedUploadContentTypes[contentType]
	if !allowed {
		http.Error(w, "tipe file tidak diizinkan", http.StatusBadRequest)
		return
	}

	kind := strings.TrimSpace(input.Kind)
	if kind == "" {
		kind = "ktp"
	}

	if kind != "ktp" && kind != "logo" {
		http.Error(w, "jenis upload tidak dikenal", http.StatusBadRequest)
		return
	}

	random, err := randomHex(16)
	if err != nil {
		log.Printf("membuat nama object: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	objectKey := fmt.Sprintf(
		"%s/%s/%s%s",
		kind,
		user.ID,
		random,
		extension,
	)

	presigned, err := app.storage.presignClient.PresignPutObject(
		r.Context(),
		&s3.PutObjectInput{
			Bucket:      aws.String(app.storage.bucket),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		},
		s3.WithPresignExpires(app.storage.presignTTL),
	)
	if err != nil {
		log.Printf("membuat presigned URL: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	if err := json.NewEncoder(w).Encode(presignUploadResponse{
		UploadURL: presigned.URL,
		ObjectKey: objectKey,
		PublicURL: app.storage.objectPublicURL(objectKey),
		ExpiresIn: int(app.storage.presignTTL.Seconds()),
	}); err != nil {
		log.Printf("encode presign response: %v", err)
	}
}
