package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims *jwt.RegisteredClaims = &jwt.RegisteredClaims{}
	var id uuid.UUID

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return id, fmt.Errorf("Error parsing token: %v", err)
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return id, fmt.Errorf("Error getting id: %v", err)
	}

	expireAt := claims.ExpiresAt.Time

	if time.Now().After(expireAt) {
		return id, fmt.Errorf("Token Expired.")
	}

	id, err = uuid.Parse(subject)
	if err != nil {
		return id, fmt.Errorf("Error parsing id: %v", err)
	}

	return id, nil
}
