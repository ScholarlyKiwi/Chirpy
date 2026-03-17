package main

import (
	"time"

	"github.com/google/uuid"
)

type jsonError struct {
	Error string `json:"error"`
}

type jsonChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type jsonLoginBody struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type jsonUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefershToken string    `json:"refresh_token"`
	IsChirpRed   bool      `json:"is_chirpy_red"`
}

type jsonToken struct {
	Token string `json:"token"`
}

type jsonNewUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type jsonWebhookData struct {
	UserID string `json:"user_id"`
}

type jsonWebhook struct {
	Event string          `json:"event"`
	Data  jsonWebhookData `json:"data"`
}
