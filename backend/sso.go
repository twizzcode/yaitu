package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type authorizeSSORequest struct {
	TargetHost string `json:"target_host"`
	ReturnPath string `json:"return_path"`
}

type authorizeSSOResponse struct {
	Code       string `json:"code"`
	TargetHost string `json:"target_host"`
	ReturnPath string `json:"return_path"`
	ExpiresIn  int    `json:"expires_in"`
}

type exchangeSSORequest struct {
	Code       string `json:"code"`
	TargetHost string `json:"target_host"`
}

type exchangeSSOResponse struct {
	SessionToken string    `json:"session_token"`
	ReturnPath   string    `json:"return_path"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func normalizeReturnPath(value string) (string, error) {
	returnPath := strings.TrimSpace(value)

	if returnPath == "" {
		returnPath = "/"
	}

	if !strings.HasPrefix(returnPath, "/") {
		return "", errors.New("return path tidak valid")
	}

	if strings.HasPrefix(returnPath, "//") {
		return "", errors.New("return path tidak valid")
	}

	if strings.ContainsAny(returnPath, "\r\n") {
		return "", errors.New("return path tidak valid")
	}

	return returnPath, nil
}

func (app *application) authorizeSSOHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.getAuthenticatedSession(r)

	if errors.Is(err, errUnauthenticated) {
		http.Error(w, "belum login", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi central session: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	if session.Hostname != app.rootDomain {
		http.Error(
			w,
			"central session diperlukan",
			http.StatusForbidden,
		)
		return
	}

	var input authorizeSSORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	targetHost, err := normalizeHostname(input.TargetHost)
	if err != nil {
		http.Error(
			w,
			"target hostname tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	returnPath, err := normalizeReturnPath(input.ReturnPath)
	if err != nil {
		http.Error(
			w,
			"return path tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	var domainIsActive bool

	err = app.db.QueryRow(
		r.Context(),
		`SELECT EXISTS (
			SELECT 1
			FROM venue_domains
			JOIN venues ON venues.id = venue_domains.venue_id
			WHERE venue_domains.hostname = $1
			  AND venue_domains.status = 'active'
			  AND venues.is_active = TRUE
		)`,
		targetHost,
	).Scan(&domainIsActive)

	if err != nil {
		log.Printf("memeriksa target SSO %q: %v", targetHost, err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	if !domainIsActive {
		http.Error(
			w,
			"target domain tidak aktif",
			http.StatusForbidden,
		)
		return
	}

	code, codeHash, err := generateSessionToken()
	if err != nil {
		log.Printf("membuat authorization code: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	expiresIn := 60
	expiresAt := time.Now().Add(
		time.Duration(expiresIn) * time.Second,
	)

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction authorize SSO: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}
	defer tx.Rollback(r.Context())

	_, err = tx.Exec(
		r.Context(),
		`DELETE FROM auth_codes
		 WHERE expires_at <= NOW()
		    OR (
		        user_id = $1
		        AND target_host = $2
		    )`,
		session.User.ID,
		targetHost,
	)
	if err != nil {
		log.Printf("membersihkan authorization code: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	_, err = tx.Exec(
		r.Context(),
		`INSERT INTO auth_codes (
			code_hash,
			user_id,
			family_id,
			target_host,
			return_path,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		codeHash,
		session.User.ID,
		session.FamilyID,
		targetHost,
		returnPath,
		expiresAt,
	)
	if err != nil {
		log.Printf("menyimpan authorization code: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit authorization code: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	if err := json.NewEncoder(w).Encode(authorizeSSOResponse{
		Code:       code,
		TargetHost: targetHost,
		ReturnPath: returnPath,
		ExpiresIn:  expiresIn,
	}); err != nil {
		log.Printf("encode authorization response: %v", err)
	}
}

func (app *application) exchangeSSOHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input exchangeSSORequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Code = strings.TrimSpace(input.Code)

	if input.Code == "" {
		http.Error(
			w,
			"authorization code wajib diisi",
			http.StatusBadRequest,
		)
		return
	}

	targetHost, err := normalizeHostname(input.TargetHost)
	if err != nil {
		http.Error(
			w,
			"target hostname tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	codeHash := sha256.Sum256([]byte(input.Code))

	sessionToken, sessionTokenHash, err := generateSessionToken()
	if err != nil {
		log.Printf("membuat tenant session token: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	sessionExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction exchange SSO: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}
	defer tx.Rollback(r.Context())

	var userID string
	var familyID string
	var returnPath string

	err = tx.QueryRow(
		r.Context(),
		`DELETE FROM auth_codes
		 USING venue_domains, venues
		 WHERE auth_codes.code_hash = $1
		   AND auth_codes.target_host = $2
		   AND auth_codes.expires_at > NOW()
		   AND venue_domains.hostname = auth_codes.target_host
		   AND venue_domains.status = 'active'
		   AND venues.id = venue_domains.venue_id
		   AND venues.is_active = TRUE
		 RETURNING
			auth_codes.user_id,
			auth_codes.family_id,
			auth_codes.return_path`,
		codeHash[:],
		targetHost,
	).Scan(
		&userID,
		&familyID,
		&returnPath,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			"authorization code tidak valid atau kedaluwarsa",
			http.StatusUnauthorized,
		)
		return
	}

	if err != nil {
		log.Printf("menukar authorization code: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	_, err = tx.Exec(
		r.Context(),
		`INSERT INTO sessions (
			token_hash,
			user_id,
			hostname,
			family_id,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (family_id, hostname)
		WHERE hostname IS NOT NULL
		DO UPDATE SET
			token_hash = EXCLUDED.token_hash,
			user_id = EXCLUDED.user_id,
			expires_at = EXCLUDED.expires_at,
			created_at = NOW()`,
		sessionTokenHash,
		userID,
		targetHost,
		familyID,
		sessionExpiresAt,
	)
	if err != nil {
		log.Printf("menyimpan tenant session: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit exchange SSO: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	if err := json.NewEncoder(w).Encode(exchangeSSOResponse{
		SessionToken: sessionToken,
		ReturnPath:   returnPath,
		ExpiresAt:    sessionExpiresAt,
	}); err != nil {
		log.Printf("encode exchange SSO: %v", err)
	}
}