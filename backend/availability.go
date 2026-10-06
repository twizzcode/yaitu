package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type availabilitySlot struct {
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	Available bool      `json:"available"`
}

type availabilityResponse struct {
	CourtID             string             `json:"court_id"`
	VenueID             string             `json:"venue_id"`
	Date                string             `json:"date"`
	Timezone            string             `json:"timezone"`
	SlotDurationMinutes int                `json:"slot_duration_minutes"`
	PricePerSlot        int                `json:"price_per_slot"`
	Slots               []availabilitySlot `json:"slots"`
}

func (app *application) courtAvailabilityHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	courtID := r.PathValue("courtID")
	if courtID == "" {
		http.Error(w, "lapangan tidak valid", http.StatusBadRequest)
		return
	}

	dateValue := r.URL.Query().Get("date")
	if dateValue == "" {
		http.Error(w, "tanggal wajib diisi", http.StatusBadRequest)
		return
	}

	requestedDate, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		http.Error(
			w,
			"format tanggal harus YYYY-MM-DD",
			http.StatusBadRequest,
		)
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
		courtID,
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
		log.Printf("mengambil lapangan untuk availability: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	location, err := time.LoadLocation(court.Timezone)
	if err != nil {
		log.Printf("memuat timezone %q: %v", court.Timezone, err)
		http.Error(w, "timezone venue tidak valid", http.StatusInternalServerError)
		return
	}

	requestedDate = time.Date(
		requestedDate.Year(),
		requestedDate.Month(),
		requestedDate.Day(),
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

	maximumDate := today.AddDate(0, 0, 30)

	if requestedDate.Before(today) {
		http.Error(
			w,
			"tanggal tidak boleh berada di masa lalu",
			http.StatusBadRequest,
		)
		return
	}

	if requestedDate.After(maximumDate) {
		http.Error(
			w,
			"tanggal maksimal 30 hari dari hari ini",
			http.StatusBadRequest,
		)
		return
	}

	dayOfWeek := (int(requestedDate.Weekday())+6)%7 + 1

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
		response := availabilityResponse{
			CourtID:             courtID,
			VenueID:             court.VenueID,
			Date:                dateValue,
			Timezone:            court.Timezone,
			SlotDurationMinutes: court.SlotDurationMinutes,
			PricePerSlot:        court.PricePerSlot,
			Slots:               make([]availabilitySlot, 0),
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("encode availability: %v", err)
		}

		return
	}

	if err != nil {
		log.Printf("mengambil jam operasional: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

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

	slotDuration := time.Duration(court.SlotDurationMinutes) * time.Minute
	slots := make([]availabilitySlot, 0)

	for startsAt := opensAt; ; startsAt = startsAt.Add(slotDuration) {
		endsAt := startsAt.Add(slotDuration)

		if endsAt.After(closesAt) {
			break
		}

		slots = append(slots, availabilitySlot{
			StartsAt:  startsAt,
			EndsAt:    endsAt,
			Available: startsAt.After(now),
		})
	}

	rows, err := app.db.Query(
		r.Context(),
		`SELECT starts_at, ends_at
		 FROM bookings
		 WHERE court_id = $1
		   AND starts_at < $2
		   AND ends_at > $3
		   AND (
		       status = 'confirmed'
		       OR (
		           status = 'pending_payment'
		           AND payment_expires_at > NOW()
		       )
		   )
		 ORDER BY starts_at`,
		courtID,
		closesAt,
		opensAt,
	)
	if err != nil {
		log.Printf("mengambil booking aktif: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	type bookedPeriod struct {
		StartsAt time.Time
		EndsAt   time.Time
	}
	
	bookedPeriods := make([]bookedPeriod, 0)
	
	for rows.Next() {
		var period bookedPeriod
	
		if err := rows.Scan(&period.StartsAt, &period.EndsAt); err != nil {
			log.Printf("membaca booking aktif: %v", err)
			http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
			return
		}
	
		bookedPeriods = append(bookedPeriods, period)
	}
	
	if err := rows.Err(); err != nil {
		log.Printf("iterasi booking aktif: %v", err)
		http.Error(w, "terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}
	
	for i := range slots {
		if !slots[i].Available {
			continue
		}
	
		for _, period := range bookedPeriods {
			hasOverlap :=
				slots[i].StartsAt.Before(period.EndsAt) &&
					slots[i].EndsAt.After(period.StartsAt)
	
			if hasOverlap {
				slots[i].Available = false
				break
			}
		}
	}

	response := availabilityResponse{
		CourtID:             courtID,
		VenueID:             court.VenueID,
		Date:                dateValue,
		Timezone:            court.Timezone,
		SlotDurationMinutes: court.SlotDurationMinutes,
		PricePerSlot:        court.PricePerSlot,
		Slots:               slots,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode availability: %v", err)
	}
}
