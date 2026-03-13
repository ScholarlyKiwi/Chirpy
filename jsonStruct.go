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
	Email    string `json:"email"`
	Password string `json:"password"`
}
