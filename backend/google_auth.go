package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

type googleAuthStartResponse struct {
	AuthorizationURL string `json:"authorization_url"`
}

type googleAuthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type googleAuthCallbackResponse struct {
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	NextPath     string    `json:"next_path"`
}

type googleIDTokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Nonce         string `json:"nonce"`
}

func generateOAuthValue() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func normalizeNextPath(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "/admin"
	}

	if !strings.HasPrefix(value, "/") ||
		strings.HasPrefix(value, "//") ||
		strings.ContainsAny(value, "\r\n") {
		return "/admin"
	}

	return value
}

func (app *application) startGoogleAuthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if app.googleOAuthConfig == nil {
		http.Error(
			w,
			"login Google tidak tersedia",
			http.StatusServiceUnavailable,
		)
		return
	}

	state, err := generateOAuthValue()
	if err != nil {
		log.Printf("membuat OAuth state: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	nonce, err := generateOAuthValue()
	if err != nil {
		log.Printf("membuat OIDC nonce: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	codeVerifier := oauth2.GenerateVerifier()
	stateHash := sha256.Sum256([]byte(state))
	nextPath := normalizeNextPath(r.URL.Query().Get("next"))
	expiresAt := time.Now().Add(10 * time.Minute)

	if _, err := app.db.Exec(
		r.Context(),
		`DELETE FROM oauth_states WHERE expires_at <= NOW()`,
	); err != nil {
		log.Printf("membersihkan OAuth state: %v", err)
	}

	if _, err := app.db.Exec(
		r.Context(),
		`INSERT INTO oauth_states (
			state_hash,
			nonce,
			code_verifier,
			next_path,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5)`,
		stateHash[:],
		nonce,
		codeVerifier,
		nextPath,
		expiresAt,
	); err != nil {
		log.Printf("menyimpan OAuth state: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	authorizationURL := app.googleOAuthConfig.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
		oauth2.SetAuthURLParam("nonce", nonce),
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		googleAuthStartResponse{
			AuthorizationURL: authorizationURL,
		},
	); err != nil {
		log.Printf("encode authorization URL: %v", err)
	}
}

func (app *application) linkGoogleUser(
	ctx context.Context,
	subject string,
	claims googleIDTokenClaims,
) (string, error) {
	tx, err := app.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("memulai transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID string

	err = tx.QueryRow(
		ctx,
		`SELECT user_id
		 FROM user_identities
		 WHERE provider = 'google'
		   AND subject = $1`,
		subject,
	).Scan(&userID)

	if err == nil {
		if _, err := tx.Exec(
			ctx,
			`UPDATE user_identities
			 SET email = $1
			 WHERE provider = 'google'
			   AND subject = $2`,
			claims.Email,
			subject,
		); err != nil {
			return "", fmt.Errorf(
				"memperbarui identitas Google: %w",
				err,
			)
		}

		if err := tx.Commit(ctx); err != nil {
			return "", fmt.Errorf(
				"commit identitas Google: %w",
				err,
			)
		}

		return userID, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf(
			"mencari identitas Google: %w",
			err,
		)
	}

	name := strings.TrimSpace(claims.Name)
	if name == "" {
		name = claims.Email
	}

	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (
			id,
			name,
			email,
			email_verified_at
		)
		VALUES (
			gen_random_uuid(),
			$1,
			$2,
			NOW()
		)
		ON CONFLICT (email) DO UPDATE
		SET
			email_verified_at = COALESCE(
				users.email_verified_at,
				EXCLUDED.email_verified_at
			),
			updated_at = NOW()
		RETURNING id`,
		name,
		claims.Email,
	).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf(
			"membuat atau menautkan user Google: %w",
			err,
		)
	}

	var linkedUserID string

	err = tx.QueryRow(
		ctx,
		`INSERT INTO user_identities (
			provider,
			subject,
			user_id,
			email
		)
		VALUES ('google', $1, $2, $3)
		ON CONFLICT (provider, subject) DO UPDATE
		SET email = EXCLUDED.email
		RETURNING user_id`,
		subject,
		userID,
		claims.Email,
	).Scan(&linkedUserID)
	if err != nil {
		return "", fmt.Errorf(
			"menyimpan identitas Google: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf(
			"commit identitas Google: %w",
			err,
		)
	}

	return linkedUserID, nil
}

func (app *application) googleAuthCallbackHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if app.googleOAuthConfig == nil || app.googleIDTokenVerifier == nil {
		http.Error(
			w,
			"login Google tidak tersedia",
			http.StatusServiceUnavailable,
		)
		return
	}

	var input googleAuthCallbackRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	input.Code = strings.TrimSpace(input.Code)
	input.State = strings.TrimSpace(input.State)

	if input.Code == "" || input.State == "" {
		http.Error(
			w,
			"authorization code dan state wajib diisi",
			http.StatusBadRequest,
		)
		return
	}

	stateHash := sha256.Sum256([]byte(input.State))

	var oauthState struct {
		Nonce        string
		CodeVerifier string
		NextPath     string
	}

	err := app.db.QueryRow(
		r.Context(),
		`DELETE FROM oauth_states
		 WHERE state_hash = $1
		   AND expires_at > NOW()
		 RETURNING nonce, code_verifier, next_path`,
		stateHash[:],
	).Scan(
		&oauthState.Nonce,
		&oauthState.CodeVerifier,
		&oauthState.NextPath,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(
			w,
			"OAuth state tidak valid atau kedaluwarsa",
			http.StatusBadRequest,
		)
		return
	}

	if err != nil {
		log.Printf("mengambil OAuth state: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	oauthToken, err := app.googleOAuthConfig.Exchange(
		r.Context(),
		input.Code,
		oauth2.VerifierOption(oauthState.CodeVerifier),
	)
	if err != nil {
		log.Printf("menukar authorization code Google: %v", err)
		http.Error(
			w,
			"authorization code Google tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		http.Error(
			w,
			"ID token Google tidak ditemukan",
			http.StatusBadRequest,
		)
		return
	}

	idToken, err := app.googleIDTokenVerifier.Verify(
		r.Context(),
		rawIDToken,
	)
	if err != nil {
		log.Printf("verifikasi ID token Google: %v", err)
		http.Error(
			w,
			"ID token Google tidak valid",
			http.StatusUnauthorized,
		)
		return
	}

	var claims googleIDTokenClaims

	if err := idToken.Claims(&claims); err != nil {
		log.Printf("membaca claim ID token Google: %v", err)
		http.Error(
			w,
			"claim Google tidak valid",
			http.StatusUnauthorized,
		)
		return
	}

	if subtle.ConstantTimeCompare(
		[]byte(claims.Nonce),
		[]byte(oauthState.Nonce),
	) != 1 {
		http.Error(
			w,
			"nonce Google tidak valid",
			http.StatusUnauthorized,
		)
		return
	}

	claims.Email = strings.ToLower(
		strings.TrimSpace(claims.Email),
	)

	if idToken.Subject == "" ||
		claims.Email == "" ||
		!claims.EmailVerified {
		http.Error(
			w,
			"akun Google tidak memiliki email terverifikasi",
			http.StatusUnauthorized,
		)
		return
	}

	userID, err := app.linkGoogleUser(
		r.Context(),
		idToken.Subject,
		claims,
	)
	if err != nil {
		log.Printf("menautkan akun Google: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	session, err := app.createCentralSession(
		r.Context(),
		userID,
	)
	if err != nil {
		log.Printf("membuat session Google: %v", err)
		http.Error(
			w,
			"terjadi kesalahan pada server",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(
		googleAuthCallbackResponse{
			SessionToken: session.Token,
			ExpiresAt:    session.ExpiresAt,
			NextPath:     oauthState.NextPath,
		},
	); err != nil {
		log.Printf("encode callback Google: %v", err)
	}
}
