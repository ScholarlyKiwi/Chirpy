package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func DecodeJWT(token string, tokenSecret string) (*jwt.Token, jwt.Claims, error) {
	var claims jwt.RegisteredClaims

	parser := &jwt.Parser{}

	decode_token, err := parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return decode_token, claims, fmt.Errorf("Error parsing token: %v", err)
	}
	return decode_token, claims, nil
}
