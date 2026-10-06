package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type authenticatedUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type authenticatedSession struct {
	User     authenticatedUser
	FamilyID string
	Hostname string
}

type createdSession struct {
	Token     string
	ExpiresAt time.Time
}

var errUnauthenticated = errors.New("user belum login")

func generateSessionToken() (string, []byte, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", nil, fmt.Errorf("membuat token session: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(randomBytes)
	tokenHash := sha256.Sum256([]byte(token))

	return token, tokenHash[:], nil
}

func (app *application) getAuthenticatedSession(r *http.Request) (authenticatedSession, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return authenticatedSession{}, errUnauthenticated
	}

	tokenHash := sha256.Sum256([]byte(cookie.Value))

	var session authenticatedSession

	err = app.db.QueryRow(
		r.Context(),
		`SELECT
			users.id,
			users.name,
			users.email,
			sessions.family_id,
			COALESCE(sessions.hostname, '')
		 FROM sessions
		 JOIN users ON users.id = sessions.user_id
		 WHERE sessions.token_hash = $1
		   AND sessions.expires_at > NOW()`,
		tokenHash[:],
	).Scan(
		&session.User.ID,
		&session.User.Name,
		&session.User.Email,
		&session.FamilyID,
		&session.Hostname,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return authenticatedSession{}, errUnauthenticated
	}

	if err != nil {
		return authenticatedSession{}, fmt.Errorf(
			"mencari session: %w",
			err,
		)
	}

	return session, nil
}

func (app *application) getAuthenticatedUser(r *http.Request) (authenticatedUser, error) {
	session, err := app.getAuthenticatedSession(r)
	if err != nil {
		return authenticatedUser{}, err
	}

	return session.User, nil
}

func (app *application) createCentralSession(ctx context.Context, userID string) (createdSession, error) {
	token, tokenHash, err := generateSessionToken()
	if err != nil {
		return createdSession{}, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	if _, err := app.db.Exec(
		ctx,
		`INSERT INTO sessions (
			token_hash,
			user_id,
			hostname,
			expires_at
		)
		VALUES ($1, $2, $3, $4)`,
		tokenHash,
		userID,
		app.rootDomain,
		expiresAt,
	); err != nil {
		return createdSession{}, fmt.Errorf(
			"menyimpan session: %w",
			err,
		)
	}

	return createdSession{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
