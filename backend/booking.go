package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type createBookingRequest struct {
	CourtID  string    `json:"court_id"`
	StartsAt time.Time `json:"starts_at"`
}

type bookingResponse struct {
	ID               string     `json:"id"`
	CourtID          string     `json:"court_id"`
	CustomerID       string     `json:"customer_id"`
	StartsAt         time.Time  `json:"starts_at"`
	EndsAt           time.Time  `json:"ends_at"`
	TotalAmount      int        `json:"total_amount"`
	Status           string     `json:"status"`
	PaymentExpiresAt *time.Time `json:"payment_expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type bookingDetailResponse struct {
	ID               string     `json:"id"`
	CourtID          string     `json:"court_id"`
	CourtName        string     `json:"court_name"`
	VenueName        string     `json:"venue_name"`
	VenueSlug        string     `json:"venue_slug"`
	Timezone         string     `json:"timezone"`
	CustomerID       string     `json:"customer_id"`
	StartsAt         time.Time  `json:"starts_at"`
	EndsAt           time.Time  `json:"ends_at"`
	TotalAmount      int        `json:"total_amount"`
	Status           string     `json:"status"`
	PaymentExpiresAt *time.Time `json:"payment_expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

func (app *application) createBookingHandler(w http.ResponseWriter, r *http.Request) {
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

	var input createBookingRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.CourtID = strings.TrimSpace(input.CourtID)

	if input.CourtID == "" {
		http.Error(w, "court_id wajib diisi", http.StatusBadRequest)
		return
	}

	if input.StartsAt.IsZero() {
		http.Error(w, "starts_at wajib diisi", http.StatusBadRequest)
		return
	}

	var court struct {
		VenueID             string
		Timezone            string
		SlotDurationMinutes int
		PricePerSlot        int
	}

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			venues.id,
			venues.timezone,
			courts.slot_duration_minutes,
			courts.price_per_slot
		 FROM courts
		 JOIN venues ON venues.id = courts.venue_id
		 WHERE courts.id::text = $1
		   AND courts.is_active = TRUE
		   AND venues.is_active = TRUE`,
		input.CourtID,
	).Scan(
		&court.VenueID,
		&court.Timezone,
		&court.SlotDurationMinutes,
		&court.PricePerSlot,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "lapangan tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil lapangan untuk booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	location, err := time.LoadLocation(court.Timezone)
	if err != nil {
		log.Printf("memuat timezone %q: %v", court.Timezone, err)
		http.Error(w, "timezone venue tidak valid", http.StatusInternalServerError)
		return
	}

	startsAt := input.StartsAt.In(location)
	endsAt := startsAt.Add(
		time.Duration(court.SlotDurationMinutes) * time.Minute,
	)

	if !startsAt.After(time.Now().In(location)) {
		http.Error(
			w,
			"waktu booking harus berada di masa depan",
			http.StatusBadRequest,
		)
		return
	}

	localDate := time.Date(
		startsAt.Year(),
		startsAt.Month(),
		startsAt.Day(),
		0,
		0,
		0,
		0,
		location,
	)

	now := time.Now().In(location)

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		location,
	)

	if localDate.After(today.AddDate(0, 0, 30)) {
		http.Error(
			w,
			"tanggal booking maksimal 30 hari dari hari ini",
			http.StatusBadRequest,
		)
		return
	}

	dayOfWeek := (int(startsAt.Weekday())+6)%7 + 1

	var operatingHour struct {
		OpensAt  string
		ClosesAt string
		IsClosed bool
	}

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			COALESCE(TO_CHAR(opens_at, 'HH24:MI'), ''),
			COALESCE(TO_CHAR(closes_at, 'HH24:MI'), ''),
			is_closed
		 FROM operating_hours
		 WHERE venue_id = $1
		   AND day_of_week = $2`,
		court.VenueID,
		dayOfWeek,
	).Scan(
		&operatingHour.OpensAt,
		&operatingHour.ClosesAt,
		&operatingHour.IsClosed,
	)

	if errors.Is(err, pgx.ErrNoRows) || operatingHour.IsClosed {
		http.Error(
			w,
			"venue tutup pada tanggal tersebut",
			http.StatusBadRequest,
		)
		return
	}

	if err != nil {
		log.Printf("mengambil jam operasional untuk booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	dateValue := startsAt.Format("2006-01-02")

	opensAt, err := time.ParseInLocation(
		"2006-01-02 15:04",
		dateValue+" "+operatingHour.OpensAt,
		location,
	)
	if err != nil {
		log.Printf("membaca jam buka %q: %v", operatingHour.OpensAt, err)
		http.Error(w, "jam operasional tidak valid", http.StatusInternalServerError)
		return
	}

	closesAt, err := time.ParseInLocation(
		"2006-01-02 15:04",
		dateValue+" "+operatingHour.ClosesAt,
		location,
	)
	if err != nil {
		log.Printf("membaca jam tutup %q: %v", operatingHour.ClosesAt, err)
		http.Error(w, "jam operasional tidak valid", http.StatusInternalServerError)
		return
	}

	if startsAt.Before(opensAt) || endsAt.After(closesAt) {
		http.Error(
			w,
			"waktu booking berada di luar jam operasional",
			http.StatusBadRequest,
		)
		return
	}

	slotDuration := time.Duration(court.SlotDurationMinutes) * time.Minute
	elapsedFromOpening := startsAt.Sub(opensAt)

	if elapsedFromOpening%slotDuration != 0 {
		http.Error(
			w,
			"waktu booking harus sesuai dengan awal slot",
			http.StatusBadRequest,
		)
		return
	}

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	_, err = tx.Exec(
		r.Context(),
		`UPDATE bookings
		 SET status = 'expired',
		     updated_at = NOW()
		 WHERE court_id = $1
		   AND status = 'pending_payment'
		   AND payment_expires_at <= NOW()`,
		input.CourtID,
	)
	if err != nil {
		log.Printf("mengakhiri booking kedaluwarsa: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	paymentExpiresAt := time.Now().Add(15 * time.Minute)

	var booking bookingResponse

	err = tx.QueryRow(
		r.Context(),
		`INSERT INTO bookings (
			court_id,
			customer_id,
			starts_at,
			ends_at,
			total_amount,
			status,
			payment_expires_at
		)
		VALUES ($1, $2, $3, $4, $5, 'pending_payment', $6)
		RETURNING
			id,
			court_id,
			customer_id,
			starts_at,
			ends_at,
			total_amount,
			status,
			payment_expires_at,
			created_at`,
		input.CourtID,
		user.ID,
		startsAt,
		endsAt,
		court.PricePerSlot,
		paymentExpiresAt,
	).Scan(
		&booking.ID,
		&booking.CourtID,
		&booking.CustomerID,
		&booking.StartsAt,
		&booking.EndsAt,
		&booking.TotalAmount,
		&booking.Status,
		&booking.PaymentExpiresAt,
		&booking.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			http.Error(
				w,
				"slot sudah dipesan atau sedang menunggu pembayaran",
				http.StatusConflict,
			)
			return
		}

		log.Printf("menyimpan booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("encode booking: %v", err)
	}
}

func (app *application) getBookingHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = app.db.Exec(
		r.Context(),
		`UPDATE bookings
		 SET status = 'expired',
		     updated_at = NOW()
		 WHERE id::text = $1
		   AND customer_id = $2
		   AND status = 'pending_payment'
		   AND payment_expires_at <= NOW()`,
		bookingID,
		user.ID,
	)
	if err != nil {
		log.Printf("memperbarui booking kedaluwarsa: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	var booking bookingDetailResponse

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			bookings.id,
			bookings.court_id,
			courts.name,
			venues.name,
			venues.slug,
			venues.timezone,
			bookings.customer_id,
			bookings.starts_at,
			bookings.ends_at,
			bookings.total_amount,
			bookings.status,
			bookings.payment_expires_at,
			bookings.created_at
		 FROM bookings
		 JOIN courts ON courts.id = bookings.court_id
		 JOIN venues ON venues.id = courts.venue_id
		 WHERE bookings.id::text = $1
		   AND bookings.customer_id = $2`,
		bookingID,
		user.ID,
	).Scan(
		&booking.ID,
		&booking.CourtID,
		&booking.CourtName,
		&booking.VenueName,
		&booking.VenueSlug,
		&booking.Timezone,
		&booking.CustomerID,
		&booking.StartsAt,
		&booking.EndsAt,
		&booking.TotalAmount,
		&booking.Status,
		&booking.PaymentExpiresAt,
		&booking.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "booking tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil detail booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("encode detail booking: %v", err)
	}
}

func (app *application) listCustomerBookingsHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = app.db.Exec(
		r.Context(),
		`UPDATE bookings
		 SET status = 'expired',
		     updated_at = NOW()
		 WHERE customer_id = $1
		   AND status = 'pending_payment'
		   AND payment_expires_at <= NOW()`,
		user.ID,
	)
	if err != nil {
		log.Printf("memperbarui booking kedaluwarsa: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	rows, err := app.db.Query(
		r.Context(),
		`SELECT
			bookings.id,
			bookings.court_id,
			courts.name,
			venues.name,
			venues.slug,
			venues.timezone,
			bookings.customer_id,
			bookings.starts_at,
			bookings.ends_at,
			bookings.total_amount,
			bookings.status,
			bookings.payment_expires_at,
			bookings.created_at
		 FROM bookings
		 JOIN courts ON courts.id = bookings.court_id
		 JOIN venues ON venues.id = courts.venue_id
		 WHERE bookings.customer_id = $1
		 ORDER BY bookings.created_at DESC`,
		user.ID,
	)
	if err != nil {
		log.Printf("mengambil daftar booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bookings := make([]bookingDetailResponse, 0)

	for rows.Next() {
		var booking bookingDetailResponse

		if err := rows.Scan(
			&booking.ID,
			&booking.CourtID,
			&booking.CourtName,
			&booking.VenueName,
			&booking.VenueSlug,
			&booking.Timezone,
			&booking.CustomerID,
			&booking.StartsAt,
			&booking.EndsAt,
			&booking.TotalAmount,
			&booking.Status,
			&booking.PaymentExpiresAt,
			&booking.CreatedAt,
		); err != nil {
			log.Printf("membaca daftar booking: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi daftar booking: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("encode daftar booking: %v", err)
	}
}