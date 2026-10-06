package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type publicCourtResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Sport               string `json:"sport"`
	PricePerSlot        int    `json:"price_per_slot"`
	SlotDurationMinutes int    `json:"slot_duration_minutes"`
}

type publicVenueResponse struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	Slug     string                `json:"slug"`
	Address  string                `json:"address"`
	Timezone string                `json:"timezone"`
	Whatsapp string                `json:"whatsapp"`
	Courts   []publicCourtResponse `json:"courts"`
}

func (app *application) publicVenueHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := strings.ToLower(strings.TrimSpace(r.PathValue("slug")))

	if slug == "" {
		http.Error(w, "slug venue tidak valid", http.StatusBadRequest)
		return
	}

	var venue publicVenueResponse

	err := app.db.QueryRow(
		r.Context(),
		`SELECT
			id,
			name,
			slug,
			address,
			timezone,
			COALESCE(whatsapp, '')
		 FROM venues
		 WHERE slug = $1
		   AND is_active = TRUE`,
		slug,
	).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Slug,
		&venue.Address,
		&venue.Timezone,
		&venue.Whatsapp,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "venue tidak ditemukan", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("mengambil venue publik: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	rows, err := app.db.Query(
		r.Context(),
		`SELECT
			id,
			name,
			sport,
			price_per_slot,
			slot_duration_minutes
		 FROM courts
		 WHERE venue_id = $1
		   AND is_active = TRUE
		 ORDER BY created_at`,
		venue.ID,
	)
	if err != nil {
		log.Printf("mengambil lapangan publik: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	venue.Courts = make([]publicCourtResponse, 0)

	for rows.Next() {
		var court publicCourtResponse

		if err := rows.Scan(
			&court.ID,
			&court.Name,
			&court.Sport,
			&court.PricePerSlot,
			&court.SlotDurationMinutes,
		); err != nil {
			log.Printf("membaca lapangan publik: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		venue.Courts = append(venue.Courts, court)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi lapangan publik: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(venue); err != nil {
		log.Printf("encode venue publik: %v", err)
	}
}