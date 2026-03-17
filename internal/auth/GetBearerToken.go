package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth_head := headers.Get("Authorization")
	if auth_head == "" {
		return "", fmt.Errorf("Unauthorized request - missing authorization header")
	}
	token := strings.Trim(strings.TrimLeft(auth_head, "Bearer"), " ")
	if token == "" {
		return "", fmt.Errorf("Unauthorized requested - no bearer token supplied %v", token)
	}
	return token, nil
}
