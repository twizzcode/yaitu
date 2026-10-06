package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

type domainResolutionResponse struct {
	Hostname  string `json:"hostname"`
	VenueID   string `json:"venue_id"`
	VenueSlug string `json:"venue_slug"`
	Type      string `json:"type"`
}
var hostnameLabelPattern = regexp.MustCompile(
	`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`,
)

var errInvalidHostname = errors.New("hostname tidak valid")

func normalizeHostname(value string) (string, error) {
	hostname := strings.ToLower(strings.TrimSpace(value))
	hostname = strings.TrimSuffix(hostname, ".")

	if hostname == "" || len(hostname) > 253 {
		return "", errInvalidHostname
	}

	if strings.Contains(hostname, "://") {
		return "", errInvalidHostname
	}

	if strings.ContainsAny(hostname, "/:#?") {
		return "", errInvalidHostname
	}

	labels := strings.Split(hostname, ".")

	if len(labels) < 2 {
		return "", errInvalidHostname
	}

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return "", errInvalidHostname
		}

		if !hostnameLabelPattern.MatchString(label) {
			return "", fmt.Errorf(
				"%w: label %q tidak valid",
				errInvalidHostname,
				label,
			)
		}
	}

	return hostname, nil
}

func (app *application) syncPlatformDomains(ctx context.Context) error {
	_, err := app.db.Exec(
		ctx,
		`INSERT INTO venue_domains (
			venue_id,
			hostname,
			type,
			status,
			verified_at
		)
		SELECT
			venues.id,
			venues.slug || '.' || $1,
			'platform',
			'active',
			NOW()
		FROM venues
		ON CONFLICT (venue_id) WHERE type = 'platform'
		DO UPDATE SET
			hostname = EXCLUDED.hostname,
			status = 'active',
			verification_token_hash = NULL,
			verified_at = NOW(),
			updated_at = NOW()`,
		app.rootDomain,
	)
	if err != nil {
		return fmt.Errorf(
			"menyinkronkan platform domain: %w",
			err,
		)
	}

	return nil
}

func (app *application) resolveDomainHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	hostname, err := normalizeHostname(
		r.URL.Query().Get("hostname"),
	)
	if err != nil {
		http.Error(
			w,
			"hostname tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	var resolution domainResolutionResponse

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			venue_domains.hostname,
			venues.id,
			venues.slug,
			venue_domains.type
		 FROM venue_domains
		 JOIN venues ON venues.id = venue_domains.venue_id
		 WHERE venue_domains.hostname = $1
		   AND venue_domains.status = 'active'
		   AND venues.is_active = TRUE`,
		hostname,
	).Scan(
		&resolution.Hostname,
		&resolution.VenueID,
		&resolution.VenueSlug,
		&resolution.Type,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			"domain tidak ditemukan",
			http.StatusNotFound,
		)
		return
	}

	if err != nil {
		log.Printf("resolve domain %q: %v", hostname, err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resolution); err != nil {
		log.Printf("encode domain resolution: %v", err)
	}
}
