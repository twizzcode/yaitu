package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

type application struct {
	db                    *pgxpool.Pool
	rootDomain            string
	serverPublicIP        string
	saasAskToken          string
	googleOAuthConfig     *oauth2.Config
	googleIDTokenVerifier *oidc.IDTokenVerifier
	storage               *storageConfig
	midtrans              *midtransConfig
}

type registerUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL belum diatur")
	}

	rootDomain := strings.ToLower(
		strings.TrimSpace(os.Getenv("ROOT_DOMAIN")),
	)

	if rootDomain == "" {
		log.Fatal("ROOT_DOMAIN belum diatur")
	}

	if strings.ContainsAny(rootDomain, "/:") {
		log.Fatal("ROOT_DOMAIN harus berupa hostname tanpa scheme atau port")
	}

	// IP publik VPS, dipakai untuk memverifikasi A record domain customer.
	serverPublicIP := strings.TrimSpace(os.Getenv("SERVER_PUBLIC_IP"))

	// Token rahasia untuk endpoint "ask" on-demand TLS Caddy.
	saasAskToken := strings.TrimSpace(os.Getenv("SAAS_ASK_TOKEN"))

	if saasAskToken == "" {
		log.Println("SAAS_ASK_TOKEN belum diatur; endpoint ask TLS dinonaktifkan")
	}

	googleClientID := strings.TrimSpace(
		os.Getenv("GOOGLE_CLIENT_ID"),
	)
	googleClientSecret := strings.TrimSpace(
		os.Getenv("GOOGLE_CLIENT_SECRET"),
	)
	googleRedirectURL := strings.TrimSpace(
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	if googleClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID belum diatur")
	}

	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET belum diatur")
	}

	if googleRedirectURL == "" {
		log.Fatal("GOOGLE_REDIRECT_URL belum diatur")
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("membuat connection pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("menghubungkan ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL")

	googleProvider, err := oidc.NewProvider(
		ctx,
		"https://accounts.google.com",
	)
	if err != nil {
		log.Fatalf("menghubungkan ke Google OIDC: %v", err)
	}

	googleOAuthConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Endpoint:     googleProvider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			oidc.ScopeEmail,
			oidc.ScopeProfile,
		},
	}

	storage, err := newStorageConfig(ctx)
	if err != nil {
		log.Fatalf("menyiapkan object storage: %v", err)
	}

	if storage == nil {
		log.Println("Object storage tidak dikonfigurasi; upload KTP dinonaktifkan")
	}

	midtrans := newMidtransConfig()

	if midtrans == nil {
		log.Println("Midtrans tidak dikonfigurasi; pembayaran QRIS dinonaktifkan")
	}

	app := &application{
		db:             db,
		rootDomain:     rootDomain,
		serverPublicIP: serverPublicIP,
		saasAskToken:   saasAskToken,
		googleOAuthConfig: googleOAuthConfig,
		googleIDTokenVerifier: googleProvider.Verifier(
			&oidc.Config{
				ClientID: googleClientID,
			},
		),
		storage: storage,
		midtrans: midtrans,
	}

	if err := app.syncPlatformDomains(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.handlerHealth)
	mux.HandleFunc("POST /api/auth/register", app.registerUserHandler)
	mux.HandleFunc("POST /api/auth/login", app.loginUserHandler)
	mux.HandleFunc("GET /api/auth/me", app.currentUserHandler)
	mux.HandleFunc("POST /api/auth/logout", app.logoutUserHandler)
	mux.HandleFunc("POST /api/venues", app.createVenueHandler)
	mux.HandleFunc("GET /api/venues", app.listVenuesHandler)
	mux.HandleFunc("POST /api/venues/{venueID}/courts", app.createCourtHandler)
	mux.HandleFunc("GET /api/venues/{venueID}/courts", app.listCourtsHandler)
	mux.HandleFunc("PUT /api/venues/{venueID}/operating-hours", app.updateOperatingHoursHandler)
	mux.HandleFunc("GET /api/venues/{venueID}/operating-hours", app.listOperatingHoursHandler)
	mux.HandleFunc("GET /api/courts/{courtID}/availability", app.courtAvailabilityHandler)
	mux.HandleFunc("POST /api/bookings", app.createBookingHandler)
	mux.HandleFunc("GET /api/bookings", app.listCustomerBookingsHandler)
	mux.HandleFunc("GET /api/bookings/{bookingID}", app.getBookingHandler)
	mux.HandleFunc("POST /api/payments", app.createPaymentHandler)
	mux.HandleFunc("GET /api/bookings/{bookingID}/payment", app.getPaymentHandler)
	mux.HandleFunc("POST /api/webhooks/midtrans", app.midtransNotificationHandler)
	mux.HandleFunc("GET /api/public/venues/{slug}", app.publicVenueHandler)
	mux.HandleFunc("GET /api/public/domains/resolve", app.resolveDomainHandler)
	mux.HandleFunc("GET /api/internal/tls/allow/{token}", app.tlsAllowHandler)
	mux.HandleFunc("POST /api/sso/authorize", app.authorizeSSOHandler)
	mux.HandleFunc("POST /api/sso/exchange", app.exchangeSSOHandler)
	mux.HandleFunc("GET /api/venues/{venueSlug}", app.getVenueBySlugHandler)
	mux.HandleFunc("PUT /api/venues/{venueSlug}", app.updateVenueProfileHandler)
	mux.HandleFunc("GET /api/venues/{venueSlug}/domains", app.listDomainsHandler)
	mux.HandleFunc("POST /api/venues/{venueSlug}/domains", app.createDomainHandler)
	mux.HandleFunc("DELETE /api/venues/{venueSlug}/domains/{domainID}", app.deleteDomainHandler)
	mux.HandleFunc("POST /api/venues/{venueSlug}/domains/{domainID}/verify", app.verifyDomainHandler)
	mux.HandleFunc("GET /api/auth/google/start", app.startGoogleAuthHandler)
	mux.HandleFunc("POST /api/auth/google/callback", app.googleAuthCallbackHandler)
	mux.HandleFunc("POST /api/uploads/presign", app.presignUploadHandler)

	log.Println("API berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (app *application) handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	}); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Name == "" {
		http.Error(w, "nama wajib diisi", http.StatusBadRequest)
		return
	}

	address, err := mail.ParseAddress(input.Email)
	if err != nil || address.Address != input.Email {
		http.Error(w, "email tidak valid", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 8 {
		http.Error(w, "password minimal 8 karakter", http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		log.Printf("hash password: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	var userID string

	err = app.db.QueryRow(
		r.Context(),
		`INSERT INTO users (id, name, email, password_hash)
		 VALUES (gen_random_uuid(), $1, $2, $3)
		 RETURNING id`,
		input.Name,
		input.Email,
		passwordHash,
	).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "email sudah terdaftar", http.StatusConflict)
			return
		}

		log.Printf("menyimpan user: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"id":      userID,
		"name":    input.Name,
		"email":   input.Email,
		"message": "registrasi berhasil",
	}); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input loginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Email == "" || input.Password == "" {
		http.Error(w, "email dan password wajib diisi", http.StatusBadRequest)
		return
	}

	var user struct {
		ID           string
		Name         string
		Email        string
		PasswordHash string
	}

	err := app.db.QueryRow(
		r.Context(),
		`SELECT id, name, email, password_hash
		 FROM users
		 WHERE email = $1
		   AND password_hash IS NOT NULL`,
		input.Email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "email atau password salah", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("mencari user: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	matched, err := verifyPassword(input.Password, user.PasswordHash)
	if err != nil {
		log.Printf("verifikasi password: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if !matched {
		http.Error(w, "email atau password salah", http.StatusUnauthorized)
		return
	}

	session, err := app.createCentralSession(
		r.Context(),
		user.ID,
	)
	if err != nil {
		log.Printf("membuat session: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  session.ExpiresAt,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]string{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"message": "login berhasil",
	}); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) currentUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.getAuthenticatedUser(r)

	if errors.Is(err, errUnauthenticated) {
		http.Error(w, "session tidak valid atau kedaluwarsa", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi user: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		tokenHash := sha256.Sum256([]byte(cookie.Value))

		if _, err := app.db.Exec(
			r.Context(),
			`DELETE FROM sessions WHERE token_hash = $1`,
			tokenHash[:],
		); err != nil {
			log.Printf("menghapus session: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
