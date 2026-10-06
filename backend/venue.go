package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var nikPattern = regexp.MustCompile(`^\d{16}$`)

var postalCodePattern = regexp.MustCompile(`^\d{5}$`)

// validObjectKey memastikan object key berada di folder milik user. Key kosong
// dianggap valid (opsional).
func (app *application) validObjectKey(objectKey string, userID string, kind string) bool {
	if objectKey == "" {
		return true
	}

	prefix := kind + "/" + userID + "/"

	return strings.HasPrefix(objectKey, prefix) &&
		!strings.Contains(objectKey, "..")
}

func (app *application) validKtpKey(ktpKey string, userID string) bool {
	return app.validObjectKey(ktpKey, userID, "ktp")
}

func (app *application) validLogoKey(logoKey string, userID string) bool {
	return app.validObjectKey(logoKey, userID, "logo")
}

type createVenueRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Address     string `json:"address"`
	Timezone    string `json:"timezone"`
	Whatsapp    string `json:"whatsapp"`
	KtpKey      string `json:"ktp_key"`
	OwnerName   string `json:"owner_name"`
	OwnerNIK    string `json:"owner_nik"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Village     string `json:"village"`
	PostalCode  string `json:"postal_code"`
	Description string `json:"description"`
	LogoKey     string `json:"logo_key"`
}

type updateVenueProfileRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Whatsapp    string `json:"whatsapp"`
	KtpKey      string `json:"ktp_key"`
	OwnerName   string `json:"owner_name"`
	OwnerNIK    string `json:"owner_nik"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Village     string `json:"village"`
	PostalCode  string `json:"postal_code"`
	Description string `json:"description"`
	LogoKey     string `json:"logo_key"`
}

type venueResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Address     string `json:"address"`
	Timezone    string `json:"timezone"`
	Whatsapp    string `json:"whatsapp"`
	IsActive    bool   `json:"is_active"`
	Role        string `json:"role"`
	KtpKey      string `json:"ktp_key"`
	KtpURL      string `json:"ktp_url"`
	OwnerName   string `json:"owner_name"`
	OwnerNIK    string `json:"owner_nik"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	Village     string `json:"village"`
	PostalCode  string `json:"postal_code"`
	Description string `json:"description"`
	LogoKey     string `json:"logo_key"`
	LogoURL     string `json:"logo_url"`
}

// venueSelectColumns adalah daftar kolom standar untuk venue beserta role
// member. Urutannya harus cocok dengan scanVenue.
const venueSelectColumns = `
	venues.id,
	venues.name,
	venues.slug,
	venues.address,
	venues.timezone,
	COALESCE(venues.whatsapp, ''),
	venues.is_active,
	venue_members.role,
	COALESCE(venues.ktp_key, ''),
	COALESCE(venues.owner_name, ''),
	COALESCE(venues.owner_nik, ''),
	COALESCE(venues.province, ''),
	COALESCE(venues.city, ''),
	COALESCE(venues.district, ''),
	COALESCE(venues.village, ''),
	COALESCE(venues.postal_code, ''),
	COALESCE(venues.description, ''),
	COALESCE(venues.logo_key, '')`

type venueScanner interface {
	Scan(dest ...any) error
}

func scanVenue(row venueScanner) (venueResponse, error) {
	var venue venueResponse

	err := row.Scan(
		&venue.ID,
		&venue.Name,
		&venue.Slug,
		&venue.Address,
		&venue.Timezone,
		&venue.Whatsapp,
		&venue.IsActive,
		&venue.Role,
		&venue.KtpKey,
		&venue.OwnerName,
		&venue.OwnerNIK,
		&venue.Province,
		&venue.City,
		&venue.District,
		&venue.Village,
		&venue.PostalCode,
		&venue.Description,
		&venue.LogoKey,
	)

	return venue, err
}

// fillVenueURLs mengisi URL publik turunan dari object key.
func (app *application) fillVenueURLs(venue *venueResponse) {
	venue.KtpURL = app.storage.objectPublicURL(venue.KtpKey)
	venue.LogoURL = app.storage.objectPublicURL(venue.LogoKey)
}

func normalizeOwnerNIK(value string) string {
	return strings.TrimSpace(value)
}

func (app *application) createVenueHandler(w http.ResponseWriter, r *http.Request) {
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

	var input createVenueRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.ToLower(
		strings.TrimSpace(input.Slug),
	)
	input.Address = strings.TrimSpace(input.Address)
	input.Timezone = strings.TrimSpace(input.Timezone)
	input.Whatsapp = strings.TrimSpace(input.Whatsapp)
	input.KtpKey = strings.TrimSpace(input.KtpKey)
	input.OwnerName = strings.TrimSpace(input.OwnerName)
	input.OwnerNIK = normalizeOwnerNIK(input.OwnerNIK)
	input.Province = strings.TrimSpace(input.Province)
	input.City = strings.TrimSpace(input.City)
	input.District = strings.TrimSpace(input.District)
	input.Village = strings.TrimSpace(input.Village)
	input.PostalCode = strings.TrimSpace(input.PostalCode)
	input.Description = strings.TrimSpace(input.Description)
	input.LogoKey = strings.TrimSpace(input.LogoKey)

	if input.Timezone == "" {
		input.Timezone = "Asia/Jakarta"
	}

	if len(input.Name) < 2 || len(input.Name) > 100 {
		http.Error(
			w,
			"nama venue harus 2 sampai 100 karakter",
			http.StatusBadRequest,
		)
		return
	}

	if len(input.Slug) < 3 || len(input.Slug) > 63 {
		http.Error(
			w,
			"slug harus 3 sampai 63 karakter",
			http.StatusBadRequest,
		)
		return
	}

	if !slugPattern.MatchString(input.Slug) {
		http.Error(
			w,
			"slug hanya boleh berisi huruf kecil, angka, dan tanda hubung",
			http.StatusBadRequest,
		)
		return
	}

	if input.Address == "" {
		http.Error(w, "alamat wajib diisi", http.StatusBadRequest)
		return
	}

	if input.OwnerNIK != "" && !nikPattern.MatchString(input.OwnerNIK) {
		http.Error(w, "NIK harus 16 digit angka", http.StatusBadRequest)
		return
	}

	if input.PostalCode != "" && !postalCodePattern.MatchString(input.PostalCode) {
		http.Error(w, "kode pos harus 5 digit angka", http.StatusBadRequest)
		return
	}

	if len(input.Description) > 500 {
		http.Error(
			w,
			"deskripsi maksimal 500 karakter",
			http.StatusBadRequest,
		)
		return
	}

	if _, err := time.LoadLocation(input.Timezone); err != nil {
		http.Error(
			w,
			"timezone tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	platformHostname, err := normalizeHostname(
		input.Slug + "." + app.rootDomain,
	)
	if err != nil {
		log.Printf("membentuk platform hostname: %v", err)
		http.Error(
			w,
			"hostname venue tidak valid",
			http.StatusInternalServerError,
		)
		return
	}

	if !app.validKtpKey(input.KtpKey, user.ID) {
		http.Error(w, "ktp_key tidak valid", http.StatusBadRequest)
		return
	}

	if !app.validLogoKey(input.LogoKey, user.ID) {
		http.Error(w, "logo_key tidak valid", http.StatusBadRequest)
		return
	}

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	venue, err := scanVenue(tx.QueryRow(
		r.Context(),
		`INSERT INTO venues (
			name,
			slug,
			address,
			timezone,
			whatsapp,
			ktp_key,
			owner_name,
			owner_nik,
			province,
			city,
			district,
			village,
			postal_code,
			description,
			logo_key
		)
		VALUES (
			$1, $2, $3, $4, $5,
			NULLIF($6, ''),
			NULLIF($7, ''),
			NULLIF($8, ''),
			NULLIF($9, ''),
			NULLIF($10, ''),
			NULLIF($11, ''),
			NULLIF($12, ''),
			NULLIF($13, ''),
			NULLIF($14, ''),
			NULLIF($15, '')
		)
		RETURNING
			id,
			name,
			slug,
			address,
			timezone,
			COALESCE(whatsapp, ''),
			is_active,
			'owner'::text,
			COALESCE(ktp_key, ''),
			COALESCE(owner_name, ''),
			COALESCE(owner_nik, ''),
			COALESCE(province, ''),
			COALESCE(city, ''),
			COALESCE(district, ''),
			COALESCE(village, ''),
			COALESCE(postal_code, ''),
			COALESCE(description, ''),
			COALESCE(logo_key, '')`,
		input.Name,
		input.Slug,
		input.Address,
		input.Timezone,
		input.Whatsapp,
		input.KtpKey,
		input.OwnerName,
		input.OwnerNIK,
		input.Province,
		input.City,
		input.District,
		input.Village,
		input.PostalCode,
		input.Description,
		input.LogoKey,
	))

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(
				w,
				"slug venue sudah digunakan",
				http.StatusConflict,
			)
			return
		}

		log.Printf("menyimpan venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	app.fillVenueURLs(&venue)

	_, err = tx.Exec(
		r.Context(),
		`INSERT INTO venue_members (venue_id, user_id, role)
		 VALUES ($1, $2, $3)`,
		venue.ID,
		user.ID,
		venue.Role,
	)
	if err != nil {
		log.Printf("menyimpan owner venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(
		r.Context(),
		`INSERT INTO venue_domains (
			venue_id,
			hostname,
			type,
			status,
			verified_at
		)
		VALUES ($1, $2, 'platform', 'active', NOW())`,
		venue.ID,
		platformHostname,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(
				w,
				"hostname venue sudah digunakan",
				http.StatusConflict,
			)
			return
		}

		log.Printf("menyimpan platform domain venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(venue); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) listVenuesHandler(w http.ResponseWriter, r *http.Request) {
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

	rows, err := app.db.Query(
		r.Context(),
		`SELECT`+venueSelectColumns+`
		 FROM venue_members
		 JOIN venues ON venues.id = venue_members.venue_id
		 WHERE venue_members.user_id = $1
		 ORDER BY venues.created_at DESC`,
		user.ID,
	)
	if err != nil {
		log.Printf("mengambil venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	venues := make([]venueResponse, 0)

	for rows.Next() {
		venue, err := scanVenue(rows)

		if err != nil {
			log.Printf("membaca venue: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		app.fillVenueURLs(&venue)

		venues = append(venues, venue)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(venues); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) getVenueBySlugHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.getAuthenticatedUser(r)
	if errors.Is(err, errUnauthenticated) {
		http.Error(w, "belum login", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("autentikasi user: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	slug := strings.ToLower(
		strings.TrimSpace(r.PathValue("venueSlug")),
	)

	if slug == "" || !slugPattern.MatchString(slug) {
		http.Error(w, "slug venue tidak valid", http.StatusBadRequest)
		return
	}

	venue, err := scanVenue(app.db.QueryRow(
		r.Context(),
		`SELECT`+venueSelectColumns+`
		 FROM venues
		 JOIN venue_members
		   ON venue_members.venue_id = venues.id
		 WHERE venues.slug = $1
		   AND venue_members.user_id = $2
		   AND venue_members.role IN ('owner', 'admin')`,
		slug,
		user.ID,
	))

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "venue tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil venue berdasarkan slug: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	app.fillVenueURLs(&venue)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(venue); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) updateVenueProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	slug := strings.ToLower(
		strings.TrimSpace(r.PathValue("venueSlug")),
	)

	if slug == "" || !slugPattern.MatchString(slug) {
		http.Error(w, "slug venue tidak valid", http.StatusBadRequest)
		return
	}

	var input updateVenueProfileRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Whatsapp = strings.TrimSpace(input.Whatsapp)
	input.KtpKey = strings.TrimSpace(input.KtpKey)
	input.OwnerName = strings.TrimSpace(input.OwnerName)
	input.OwnerNIK = normalizeOwnerNIK(input.OwnerNIK)
	input.Province = strings.TrimSpace(input.Province)
	input.City = strings.TrimSpace(input.City)
	input.District = strings.TrimSpace(input.District)
	input.Village = strings.TrimSpace(input.Village)
	input.PostalCode = strings.TrimSpace(input.PostalCode)
	input.Description = strings.TrimSpace(input.Description)
	input.LogoKey = strings.TrimSpace(input.LogoKey)

	if len(input.Name) < 2 || len(input.Name) > 100 {
		http.Error(
			w,
			"nama venue harus 2 sampai 100 karakter",
			http.StatusBadRequest,
		)
		return
	}

	if input.Address == "" {
		http.Error(w, "alamat wajib diisi", http.StatusBadRequest)
		return
	}

	if input.OwnerNIK != "" && !nikPattern.MatchString(input.OwnerNIK) {
		http.Error(w, "NIK harus 16 digit angka", http.StatusBadRequest)
		return
	}

	if input.PostalCode != "" && !postalCodePattern.MatchString(input.PostalCode) {
		http.Error(w, "kode pos harus 5 digit angka", http.StatusBadRequest)
		return
	}

	if len(input.Description) > 500 {
		http.Error(
			w,
			"deskripsi maksimal 500 karakter",
			http.StatusBadRequest,
		)
		return
	}

	if !app.validKtpKey(input.KtpKey, user.ID) {
		http.Error(w, "ktp_key tidak valid", http.StatusBadRequest)
		return
	}

	if !app.validLogoKey(input.LogoKey, user.ID) {
		http.Error(w, "logo_key tidak valid", http.StatusBadRequest)
		return
	}

	venue, err := scanVenue(app.db.QueryRow(
		r.Context(),
		`UPDATE venues
		 SET
			name = $3,
			address = $4,
			whatsapp = NULLIF($5, ''),
			ktp_key = COALESCE(NULLIF($6, ''), ktp_key),
			owner_name = NULLIF($7, ''),
			owner_nik = NULLIF($8, ''),
			province = NULLIF($9, ''),
			city = NULLIF($10, ''),
			district = NULLIF($11, ''),
			village = NULLIF($12, ''),
			postal_code = NULLIF($13, ''),
			description = NULLIF($14, ''),
			logo_key = COALESCE(NULLIF($15, ''), logo_key),
			updated_at = NOW()
		 FROM venue_members
		 WHERE venues.id = venue_members.venue_id
		   AND venues.slug = $1
		   AND venue_members.user_id = $2
		   AND venue_members.role IN ('owner', 'admin')
		 RETURNING`+venueSelectColumns,
		slug,
		user.ID,
		input.Name,
		input.Address,
		input.Whatsapp,
		input.KtpKey,
		input.OwnerName,
		input.OwnerNIK,
		input.Province,
		input.City,
		input.District,
		input.Village,
		input.PostalCode,
		input.Description,
		input.LogoKey,
	))

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "venue tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("memperbarui profil venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	app.fillVenueURLs(&venue)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(venue); err != nil {
		log.Printf("encode response: %v", err)
	}
}
