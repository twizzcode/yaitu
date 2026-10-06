package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type domainResponse struct {
	ID                string     `json:"id"`
	Hostname          string     `json:"hostname"`
	Type              string     `json:"type"`
	Status            string     `json:"status"`
	VerificationToken string     `json:"verification_token"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CreatedAt         time.Time  `json:"created_at"`
	// IsApex true bila hostname adalah apex (tanpa subdomain), mis. lapangan-a.com.
	IsApex bool `json:"is_apex"`
}

type createDomainRequest struct {
	Hostname string `json:"hostname"`
}

const domainSelectColumns = `
	venue_domains.id,
	venue_domains.hostname,
	venue_domains.type,
	venue_domains.status,
	COALESCE(venue_domains.verification_token, ''),
	venue_domains.verified_at,
	venue_domains.created_at`

func scanDomain(row venueScanner) (domainResponse, error) {
	var domain domainResponse

	err := row.Scan(
		&domain.ID,
		&domain.Hostname,
		&domain.Type,
		&domain.Status,
		&domain.VerificationToken,
		&domain.VerifiedAt,
		&domain.CreatedAt,
	)

	if err == nil {
		domain.IsApex = isApexHostname(domain.Hostname)
	}

	return domain, err
}

// isApexHostname mengembalikan true bila hostname hanya punya dua label
// (mis. "lapangan-a.com"), yang berarti apex/root domain.
func isApexHostname(hostname string) bool {
	labels := strings.Split(hostname, ".")
	return len(labels) == 2
}

// requireVenueAccess mengembalikan venue ID bila user punya akses ke venue
// dengan slug tersebut. Menulis respons error bila tidak punya akses.
func (app *application) requireVenueAccess(
	w http.ResponseWriter,
	r *http.Request,
	userID string,
	slug string,
) (string, bool) {
	if slug == "" || !slugPattern.MatchString(slug) {
		http.Error(w, "slug venue tidak valid", http.StatusBadRequest)
		return "", false
	}

	var venueID string

	err := app.db.QueryRow(
		r.Context(),
		`SELECT venues.id
		 FROM venues
		 JOIN venue_members
		   ON venue_members.venue_id = venues.id
		 WHERE venues.slug = $1
		   AND venue_members.user_id = $2
		   AND venue_members.role IN ('owner', 'admin')`,
		slug,
		userID,
	).Scan(&venueID)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "venue tidak ditemukan", http.StatusNotFound)
		return "", false
	}

	if err != nil {
		log.Printf("memeriksa akses venue %q: %v", slug, err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return "", false
	}

	return venueID, true
}

func (app *application) listDomainsHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := strings.ToLower(strings.TrimSpace(r.PathValue("venueSlug")))

	venueID, ok := app.requireVenueAccess(w, r, user.ID, slug)
	if !ok {
		return
	}

	rows, err := app.db.Query(
		r.Context(),
		`SELECT`+domainSelectColumns+`
		 FROM venue_domains
		 WHERE venue_domains.venue_id = $1
		 ORDER BY
		   venue_domains.type = 'platform' DESC,
		   venue_domains.created_at DESC`,
		venueID,
	)
	if err != nil {
		log.Printf("mengambil domain: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	domains := make([]domainResponse, 0)

	for rows.Next() {
		domain, err := scanDomain(rows)

		if err != nil {
			log.Printf("membaca domain: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		domains = append(domains, domain)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi domain: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(domains); err != nil {
		log.Printf("encode domain: %v", err)
	}
}

func (app *application) createDomainHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := strings.ToLower(strings.TrimSpace(r.PathValue("venueSlug")))

	venueID, ok := app.requireVenueAccess(w, r, user.ID, slug)
	if !ok {
		return
	}

	var input createDomainRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	hostname, err := normalizeHostname(input.Hostname)
	if err != nil {
		http.Error(w, "hostname tidak valid", http.StatusBadRequest)
		return
	}

	// Jangan izinkan menambah domain platform milik root domain lewat sini.
	if hostname == app.rootDomain ||
		strings.HasSuffix(hostname, "."+app.rootDomain) {
		http.Error(
			w,
			"hostname berada di bawah domain platform",
			http.StatusBadRequest,
		)
		return
	}

	token, err := generateOAuthValue()
	if err != nil {
		log.Printf("membuat token verifikasi: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	domain, err := scanDomain(app.db.QueryRow(
		r.Context(),
		`INSERT INTO venue_domains (
			venue_id,
			hostname,
			type,
			status,
			verification_token
		)
		VALUES ($1, $2, 'custom', 'pending', $3)
		RETURNING`+domainSelectColumns,
		venueID,
		hostname,
		token,
	))

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(
				w,
				"domain sudah digunakan",
				http.StatusConflict,
			)
			return
		}

		log.Printf("menyimpan domain: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(domain); err != nil {
		log.Printf("encode domain: %v", err)
	}
}

func (app *application) deleteDomainHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := strings.ToLower(strings.TrimSpace(r.PathValue("venueSlug")))

	venueID, ok := app.requireVenueAccess(w, r, user.ID, slug)
	if !ok {
		return
	}

	domainID := strings.TrimSpace(r.PathValue("domainID"))
	if domainID == "" {
		http.Error(w, "domain tidak valid", http.StatusBadRequest)
		return
	}

	tag, err := app.db.Exec(
		r.Context(),
		`DELETE FROM venue_domains
		 WHERE id::text = $1
		   AND venue_id = $2
		   AND type = 'custom'`,
		domainID,
		venueID,
	)
	if err != nil {
		log.Printf("menghapus domain: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		http.Error(w, "domain tidak ditemukan", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// lookupHost mengembalikan daftar IP (A record) untuk hostname. Bisa diganti
// pada test.
var lookupHost = func(ctx context.Context, hostname string) ([]string, error) {
	return net.DefaultResolver.LookupHost(ctx, hostname)
}

// domainPointsToServer mengecek apakah A record hostname mengarah ke IP VPS.
func (app *application) domainPointsToServer(
	ctx context.Context,
	hostname string,
) bool {
	if app.serverPublicIP == "" {
		return false
	}

	ips, err := lookupHost(ctx, hostname)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if ip == app.serverPublicIP {
			return true
		}
	}

	return false
}

func (app *application) verifyDomainHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := strings.ToLower(strings.TrimSpace(r.PathValue("venueSlug")))

	venueID, ok := app.requireVenueAccess(w, r, user.ID, slug)
	if !ok {
		return
	}

	domainID := strings.TrimSpace(r.PathValue("domainID"))
	if domainID == "" {
		http.Error(w, "domain tidak valid", http.StatusBadRequest)
		return
	}

	var domain domainResponse

	domain, err = scanDomain(app.db.QueryRow(
		r.Context(),
		`SELECT`+domainSelectColumns+`
		 FROM venue_domains
		 WHERE venue_domains.id::text = $1
		   AND venue_domains.venue_id = $2
		   AND venue_domains.type = 'custom'`,
		domainID,
		venueID,
	))

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "domain tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil domain untuk verifikasi: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if domain.Status == "active" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(domain)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	verified := app.domainPointsToServer(ctx, domain.Hostname)

	newStatus := "failed"
	var verifiedAt *time.Time

	if verified {
		newStatus = "active"
		now := time.Now()
		verifiedAt = &now
	}

	domain, err = scanDomain(app.db.QueryRow(
		r.Context(),
		`UPDATE venue_domains
		 SET status = $3,
		     verified_at = $4,
		     updated_at = NOW()
		 WHERE id::text = $1
		   AND venue_id = $2
		 RETURNING`+domainSelectColumns,
		domainID,
		venueID,
		newStatus,
		verifiedAt,
	))

	if err != nil {
		log.Printf("memperbarui status domain: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !verified {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	if err := json.NewEncoder(w).Encode(domain); err != nil {
		log.Printf("encode domain: %v", err)
	}
}

// tlsAllowHandler adalah endpoint "ask" untuk on-demand TLS Caddy.
//
// Caddy memanggil endpoint ini sebelum menerbitkan sertifikat untuk sebuah
// hostname. Kita hanya mengizinkan hostname yang terdaftar di venue_domains
// (status pending atau active), sehingga orang lain tidak bisa memaksa kita
// menerbitkan sertifikat untuk domain sembarangan.
func (app *application) tlsAllowHandler(w http.ResponseWriter, r *http.Request) {
	if app.saasAskToken == "" {
		http.Error(w, "fitur tidak aktif", http.StatusServiceUnavailable)
		return
	}

	token := strings.TrimSpace(r.PathValue("token"))

	if token != app.saasAskToken {
		http.Error(w, "tidak diizinkan", http.StatusForbidden)
		return
	}

	hostname := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("domain")))
	if hostname == "" {
		http.Error(w, "domain wajib diisi", http.StatusBadRequest)
		return
	}

	var allowed bool

	err := app.db.QueryRow(
		r.Context(),
		`SELECT EXISTS (
			SELECT 1
			FROM venue_domains
			WHERE hostname = $1
			  AND type = 'custom'
			  AND status IN ('pending', 'active')
		)`,
		hostname,
	).Scan(&allowed)

	if err != nil {
		log.Printf("memeriksa izin TLS %q: %v", hostname, err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "domain tidak terdaftar", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}
