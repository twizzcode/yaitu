package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type createCourtRequest struct {
	Name                string `json:"name"`
	Sport               string `json:"sport"`
	PricePerSlot        int    `json:"price_per_slot"`
	SlotDurationMinutes int    `json:"slot_duration_minutes"`
}

type courtResponse struct {
	ID                  string `json:"id"`
	VenueID             string `json:"venue_id"`
	Name                string `json:"name"`
	Sport               string `json:"sport"`
	PricePerSlot        int    `json:"price_per_slot"`
	SlotDurationMinutes int    `json:"slot_duration_minutes"`
	IsActive            bool   `json:"is_active"`
}

func (app *application) hasVenueAccess(
	ctx context.Context,
	venueID string,
	userID string,
) (bool, error) {
	var allowed bool

	err := app.db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM venues
			JOIN venue_members
			  ON venue_members.venue_id = venues.id
			WHERE venues.id::text = $1
			  AND venue_members.user_id = $2
			  AND venue_members.role IN ('owner', 'admin')
		)`,
		venueID,
		userID,
	).Scan(&allowed)

	return allowed, err
}

func (app *application) createCourtHandler(w http.ResponseWriter, r *http.Request) {
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

	var input createCourtRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Sport = strings.ToLower(strings.TrimSpace(input.Sport))

	if len(input.Name) < 2 || len(input.Name) > 100 {
		http.Error(w, "nama lapangan harus 2 sampai 100 karakter", http.StatusBadRequest)
		return
	}

	if len(input.Sport) < 2 || len(input.Sport) > 50 {
		http.Error(w, "jenis olahraga harus 2 sampai 50 karakter", http.StatusBadRequest)
		return
	}

	if input.PricePerSlot < 0 {
		http.Error(w, "harga tidak boleh negatif", http.StatusBadRequest)
		return
	}

	if input.SlotDurationMinutes == 0 {
		input.SlotDurationMinutes = 60
	}

	if input.SlotDurationMinutes < 15 || input.SlotDurationMinutes > 1440 {
		http.Error(
			w,
			"durasi slot harus 15 sampai 1440 menit",
			http.StatusBadRequest,
		)
		return
	}

	var court courtResponse

	err = app.db.QueryRow(
		r.Context(),
		`INSERT INTO courts (
			venue_id,
			name,
			sport,
			price_per_slot,
			slot_duration_minutes
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			venue_id,
			name,
			sport,
			price_per_slot,
			slot_duration_minutes,
			is_active`,
		venueID,
		input.Name,
		input.Sport,
		input.PricePerSlot,
		input.SlotDurationMinutes,
	).Scan(
		&court.ID,
		&court.VenueID,
		&court.Name,
		&court.Sport,
		&court.PricePerSlot,
		&court.SlotDurationMinutes,
		&court.IsActive,
	)

	if err != nil {
		log.Printf("menyimpan lapangan: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(court); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func (app *application) listCourtsHandler(w http.ResponseWriter, r *http.Request) {
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
			id,
			venue_id,
			name,
			sport,
			price_per_slot,
			slot_duration_minutes,
			is_active
		 FROM courts
		 WHERE venue_id = $1
		 ORDER BY created_at DESC`,
		venueID,
	)
	if err != nil {
		log.Printf("mengambil lapangan: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	courts := make([]courtResponse, 0)

	for rows.Next() {
		var court courtResponse

		if err := rows.Scan(
			&court.ID,
			&court.VenueID,
			&court.Name,
			&court.Sport,
			&court.PricePerSlot,
			&court.SlotDurationMinutes,
			&court.IsActive,
		); err != nil {
			log.Printf("membaca lapangan: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}

		courts = append(courts, court)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterasi lapangan: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(courts); err != nil {
		log.Printf("encode response: %v", err)
	}
}
