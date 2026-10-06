package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type operatingHour struct {
	DayOfWeek int     `json:"day_of_week"`
	OpensAt   *string `json:"opens_at"`
	ClosesAt  *string `json:"closes_at"`
	IsClosed  bool    `json:"is_closed"`
}

type updateOperatingHoursRequest struct {
	Hours []operatingHour `json:"hours"`
}

func (app *application) updateOperatingHoursHandler(w http.ResponseWriter, r *http.Request) {
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

	venueID := r.PathValue("venueID")
	if venueID == "" {
		http.Error(w, "venue tidak valid", http.StatusBadRequest)
		return
	}

	allowed, err := app.hasVenueAccess(r.Context(), venueID, user.ID)
	if err != nil {
		log.Printf("memeriksa akses venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "tidak punya akses ke venue", http.StatusForbidden)
		return
	}

	var input updateOperatingHoursRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	if len(input.Hours) != 7 {
		http.Error(w, "jam operasional harus berisi 7 hari", http.StatusBadRequest)
		return
	}

	seenDays := make(map[int]bool, 7)

	for i := range input.Hours {
		hour := &input.Hours[i]

		if hour.DayOfWeek < 1 || hour.DayOfWeek > 7 {
			http.Error(w, "hari harus bernilai 1 sampai 7", http.StatusBadRequest)
			return
		}

		if seenDays[hour.DayOfWeek] {
			http.Error(w, "hari tidak boleh duplikat", http.StatusBadRequest)
			return
		}

		seenDays[hour.DayOfWeek] = true

		if hour.IsClosed {
			if hour.OpensAt != nil || hour.ClosesAt != nil {
				http.Error(
					w,
					"hari tutup tidak boleh memiliki jam buka atau tutup",
					http.StatusBadRequest,
				)
				return
			}

			continue
		}

		if hour.OpensAt == nil || hour.ClosesAt == nil {
			http.Error(
				w,
				"hari buka wajib memiliki jam buka dan tutup",
				http.StatusBadRequest,
			)
			return
		}

		opensAt := strings.TrimSpace(*hour.OpensAt)
		closesAt := strings.TrimSpace(*hour.ClosesAt)

		parsedOpensAt, err := time.Parse("15:04", opensAt)
		if err != nil {
			http.Error(w, "format jam buka harus HH:MM", http.StatusBadRequest)
			return
		}

		parsedClosesAt, err := time.Parse("15:04", closesAt)
		if err != nil {
			http.Error(w, "format jam tutup harus HH:MM", http.StatusBadRequest)
			return
		}

		if !parsedOpensAt.Before(parsedClosesAt) {
			http.Error(
				w,
				"jam buka harus sebelum jam tutup",
				http.StatusBadRequest,
			)
			return
		}

		normalizedOpensAt := parsedOpensAt.Format("15:04")
		normalizedClosesAt := parsedClosesAt.Format("15:04")

		hour.OpensAt = &normalizedOpensAt
		hour.ClosesAt = &normalizedClosesAt
	}

	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("memulai transaction jam operasional: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	for _, hour := range input.Hours {
		_, err := tx.Exec(
			r.Context(),
			`INSERT INTO operating_hours (
				venue_id,
				day_of_week,
				opens_at,
				closes_at,
				is_closed
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (venue_id, day_of_week)
			DO UPDATE SET
				opens_at = EXCLUDED.opens_at,
				closes_at = EXCLUDED.closes_at,
				is_closed = EXCLUDED.is_closed`,
			venueID,
			hour.DayOfWeek,
			hour.OpensAt,
			hour.ClosesAt,
			hour.IsClosed,
		)
		if err != nil {
			log.Printf("menyimpan jam operasional: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Printf("commit jam operasional: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) listOperatingHoursHandler(w http.ResponseWriter, r *http.Request) {
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

	venueID := r.PathValue("venueID")
	if venueID == "" {
		http.Error(w, "venue tidak valid", http.StatusBadRequest)
		return
	}

	allowed, err := app.hasVenueAccess(r.Context(), venueID, user.ID)
	if err != nil {
		log.Printf("memeriksa akses venue: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "tidak punya akses ke venue", http.StatusForbidden)
		return
	}

	rows, err := app.db.Query(
		r.Context(),
		`SELECT
			day_of_week,
			COALESCE(TO_CHAR(opens_at, 'HH24:MI'), ''),
			COALESCE(TO_CHAR(closes_at, 'HH24:MI'), ''),
			is_closed
		 FROM operating_hours
		 WHERE venue_id = $1
		 ORDER BY day_of_week`,
		venueID,
	)
	if err != nil {
		log.Printf("mengambil jam operasional: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	hours := make([]operatingHour, 0, 7)

	for rows.Next() {
		var hour operatingHour
		var opensAt string
		var closesAt string

		if err := rows.Scan(
			&hour.DayOfWeek,
			&opensAt,
			&closesAt,
			&hour.IsClosed,
		); err != nil {
			log.Printf("membaca jam operasional: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		if !hour.IsClosed {
			hour.OpensAt = &opensAt
			hour.ClosesAt = &closesAt
		}

		hours = append(hours, hour)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi jam operasional: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string][]operatingHour{
		"hours": hours,
	}); err != nil {
		log.Printf("encode response: %v", err)
	}
}
