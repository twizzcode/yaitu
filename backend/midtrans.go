package main

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// midtransConfig menyimpan konfigurasi Midtrans Core API.
type midtransConfig struct {
	serverKey    string
	clientKey    string
	isProduction bool
	baseURL      string
}

func newMidtransConfig() *midtransConfig {
	serverKey := strings.TrimSpace(os.Getenv("MIDTRANS_SERVER_KEY"))
	clientKey := strings.TrimSpace(os.Getenv("MIDTRANS_CLIENT_KEY"))
	isProduction := strings.TrimSpace(os.Getenv("MIDTRANS_IS_PRODUCTION")) == "true"

	if serverKey == "" {
		return nil
	}

	baseURL := "https://api.sandbox.midtrans.com"
	if isProduction {
		baseURL = "https://api.midtrans.com"
	}

	return &midtransConfig{
		serverKey:    serverKey,
		clientKey:    clientKey,
		isProduction: isProduction,
		baseURL:      baseURL,
	}
}

func (m *midtransConfig) authHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(m.serverKey+":"),
	)
}

// paymentResponse dikirim ke frontend.
type paymentResponse struct {
	ID            string     `json:"id"`
	BookingID     string     `json:"booking_id"`
	OrderID       string     `json:"order_id"`
	Status        string     `json:"status"`
	GrossAmount   int        `json:"gross_amount"`
	QRURL         string     `json:"qr_url"`
	PaymentType   string     `json:"payment_type"`
	ExpiresAt     *time.Time `json:"expires_at"`
	TransactionID string     `json:"transaction_id"`
}

type createPaymentRequest struct {
	BookingID string `json:"booking_id"`
}

// midtransChargeResponse adalah subset respons /v2/charge.
type midtransChargeResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionStatus string `json:"transaction_status"`
	Actions           []struct {
		Name   string `json:"name"`
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"actions"`
}

func (m *midtransConfig) chargeQRIS(
	ctx context.Context,
	orderID string,
	grossAmount int,
) (*midtransChargeResponse, error) {
	payload := map[string]any{
		"payment_type": "qris",
		"transaction_details": map[string]any{
			"order_id":     orderID,
			"gross_amount": grossAmount,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload charge: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		m.baseURL+"/v2/charge",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("membuat request charge: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", m.authHeader())

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mengirim charge: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("membaca respons charge: %w", err)
	}

	var charge midtransChargeResponse

	if err := json.Unmarshal(raw, &charge); err != nil {
		return nil, fmt.Errorf("decode respons charge: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"midtrans charge gagal (%d): %s",
			resp.StatusCode,
			charge.StatusMessage,
		)
	}

	return &charge, nil
}

// qrURLFromCharge mengambil URL gambar QR dari actions charge.
func qrURLFromCharge(charge *midtransChargeResponse) string {
	for _, action := range charge.Actions {
		if action.Name == "generate-qr-code" {
			return action.URL
		}
	}

	return ""
}

// notificationPayload adalah body webhook Midtrans.
type notificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
}

// verifySignature memverifikasi signature_key Midtrans:
// sha512(order_id + status_code + gross_amount + server_key).
func (m *midtransConfig) verifySignature(payload *notificationPayload) bool {
	input := payload.OrderID +
		payload.StatusCode +
		payload.GrossAmount +
		m.serverKey

	sum := sha512.Sum512([]byte(input))
	expected := hex.EncodeToString(sum[:])

	return expected == payload.SignatureKey
}

// mapMidtransStatus memetakan status Midtrans ke status pembayaran internal.
func mapMidtransStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return "settlement"
		}
		return "capture"
	case "settlement":
		return "settlement"
	case "pending":
		return "pending"
	case "deny":
		return "deny"
	case "cancel":
		return "cancel"
	case "expire":
		return "expire"
	case "refund", "partial_refund":
		return "refund"
	default:
		return "pending"
	}
}

func (app *application) createPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if app.midtrans == nil {
		http.Error(w, "pembayaran belum dikonfigurasi", http.StatusServiceUnavailable)
		return
	}

	user, err := app.getAuthenticatedUser(r)
	if errors.Is(err, errUnauthenticated) {
		http.Error(w, "belum login", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi user: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	var input createPaymentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.BookingID = strings.TrimSpace(input.BookingID)

	if input.BookingID == "" {
		http.Error(w, "booking_id wajib diisi", http.StatusBadRequest)
		return
	}

	var booking struct {
		ID              string
		TotalAmount     int
		Status          string
		PaymentExpiresAt *time.Time
	}

	err = app.db.QueryRow(
		r.Context(),
		`SELECT id, total_amount, status, payment_expires_at
		 FROM bookings
		 WHERE id::text = $1
		   AND customer_id = $2`,
		input.BookingID,
		user.ID,
	).Scan(
		&booking.ID,
		&booking.TotalAmount,
		&booking.Status,
		&booking.PaymentExpiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "booking tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil booking untuk pembayaran: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if booking.Status != "pending_payment" {
		http.Error(w, "booking tidak menunggu pembayaran", http.StatusConflict)
		return
	}

	if booking.PaymentExpiresAt != nil && booking.PaymentExpiresAt.Before(time.Now()) {
		http.Error(w, "waktu pembayaran sudah habis", http.StatusConflict)
		return
	}

	// Bila sudah ada payment pending untuk booking ini, kembalikan yang ada
	// agar tidak membuat transaksi Midtrans baru (idempotent).
	var existing paymentResponse

	err = app.db.QueryRow(
		r.Context(),
		`SELECT id, booking_id, order_id, status, gross_amount,
		        COALESCE(qr_url, ''), COALESCE(payment_type, ''),
		        expires_at, COALESCE(transaction_id, '')
		 FROM payments
		 WHERE booking_id = $1
		   AND status = 'pending'
		 ORDER BY created_at DESC
		 LIMIT 1`,
		booking.ID,
	).Scan(
		&existing.ID,
		&existing.BookingID,
		&existing.OrderID,
		&existing.Status,
		&existing.GrossAmount,
		&existing.QRURL,
		&existing.PaymentType,
		&existing.ExpiresAt,
		&existing.TransactionID,
	)

	if err == nil && existing.QRURL != "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(existing)
		return
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("mencari payment existing: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	orderID := fmt.Sprintf("BK-%s-%d", booking.ID, time.Now().Unix())

	charge, err := app.midtrans.chargeQRIS(
		r.Context(),
		orderID,
		booking.TotalAmount,
	)
	if err != nil {
		log.Printf("charge QRIS: %v", err)
		http.Error(w, "gagal membuat pembayaran QRIS", http.StatusBadGateway)
		return
	}

	qrURL := qrURLFromCharge(charge)

	rawCharge, _ := json.Marshal(charge)

	var payment paymentResponse

	err = app.db.QueryRow(
		r.Context(),
		`INSERT INTO payments (
			booking_id,
			user_id,
			order_id,
			transaction_id,
			payment_type,
			gross_amount,
			status,
			qr_url,
			expires_at,
			raw_charge
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending', $7, $8, $9)
		RETURNING
			id, booking_id, order_id, status, gross_amount,
			COALESCE(qr_url, ''), COALESCE(payment_type, ''),
			expires_at, COALESCE(transaction_id, '')`,
		booking.ID,
		user.ID,
		orderID,
		charge.TransactionID,
		charge.PaymentType,
		booking.TotalAmount,
		qrURL,
		booking.PaymentExpiresAt,
		rawCharge,
	).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Status,
		&payment.GrossAmount,
		&payment.QRURL,
		&payment.PaymentType,
		&payment.ExpiresAt,
		&payment.TransactionID,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "pembayaran sudah dibuat", http.StatusConflict)
			return
		}

		log.Printf("menyimpan payment: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(payment); err != nil {
		log.Printf("encode payment: %v", err)
	}
}

// getPaymentHandler mengembalikan payment terbaru untuk sebuah booking.
func (app *application) getPaymentHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.getAuthenticatedUser(r)
	if errors.Is(err, errUnauthenticated) {
		http.Error(w, "belum login", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi user: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	bookingID := strings.TrimSpace(r.PathValue("bookingID"))
	if bookingID == "" {
		http.Error(w, "booking tidak valid", http.StatusBadRequest)
		return
	}

	var payment paymentResponse

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			payments.id,
			payments.booking_id,
			payments.order_id,
			payments.status,
			payments.gross_amount,
			COALESCE(payments.qr_url, ''),
			COALESCE(payments.payment_type, ''),
			payments.expires_at,
			COALESCE(payments.transaction_id, '')
		 FROM payments
		 JOIN bookings ON bookings.id = payments.booking_id
		 WHERE payments.booking_id::text = $1
		   AND bookings.customer_id = $2
		 ORDER BY payments.created_at DESC
		 LIMIT 1`,
		bookingID,
		user.ID,
	).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Status,
		&payment.GrossAmount,
		&payment.QRURL,
		&payment.PaymentType,
		&payment.ExpiresAt,
		&payment.TransactionID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "pembayaran tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil payment: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(payment); err != nil {
		log.Printf("encode payment: %v", err)
	}
}

// midtransNotificationHandler menerima webhook dari Midtrans.
func (app *application) midtransNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if app.midtrans == nil {
		http.Error(w, "pembayaran belum dikonfigurasi", http.StatusServiceUnavailable)
		return
	}

	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "body tidak valid", http.StatusBadRequest)
		return
	}

	var payload notificationPayload

	if err := json.Unmarshal(raw, &payload); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	if !app.midtrans.verifySignature(&payload) {
		log.Printf("signature webhook Midtrans tidak valid untuk order %q", payload.OrderID)
		http.Error(w, "signature tidak valid", http.StatusUnauthorized)
		return
	}

	status := mapMidtransStatus(payload.TransactionStatus, payload.FraudStatus)

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction webhook: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Ambil booking_id dan pastikan order cocok.
	var paymentID string
	var bookingID *string

	err = tx.QueryRow(
		r.Context(),
		`UPDATE payments
		 SET status = $2,
		     transaction_id = COALESCE(NULLIF($3, ''), transaction_id),
		     raw_notification = $4,
		     updated_at = NOW()
		 WHERE order_id = $1
		 RETURNING id, booking_id`,
		payload.OrderID,
		status,
		payload.TransactionID,
		raw,
	).Scan(&paymentID, &bookingID)

	if errors.Is(err, pgx.ErrNoRows) {
		// Order tidak dikenal; balas 200 agar Midtrans tidak retry tanpa henti.
		log.Printf("webhook untuk order tidak dikenal: %q", payload.OrderID)
		w.WriteHeader(http.StatusOK)
		return
	}

	if err != nil {
		log.Printf("memperbarui payment dari webhook: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	// Update status booking bila pembayaran sukses/gagal.
	if bookingID != nil {
		switch status {
		case "settlement", "capture":
			if _, err := tx.Exec(
				r.Context(),
				`UPDATE bookings
				 SET status = 'confirmed',
				     updated_at = NOW()
				 WHERE id = $1
				   AND status = 'pending_payment'`,
				*bookingID,
			); err != nil {
				log.Printf("mengonfirmasi booking: %v", err)
				http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
				return
			}
		case "expire", "cancel", "deny", "failure":
			if _, err := tx.Exec(
				r.Context(),
				`UPDATE bookings
				 SET status = 'expired',
				     updated_at = NOW()
				 WHERE id = $1
				   AND status = 'pending_payment'`,
				*bookingID,
			); err != nil {
				log.Printf("mengekspirasi booking: %v", err)
				http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit webhook: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
