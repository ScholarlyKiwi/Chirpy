package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth_head := headers.Get("Authorization")
	if auth_head == "" {
		return "", fmt.Errorf("Unauthorized request - missing authorization header")
	}

	apiKey := strings.Trim(strings.TrimLeft(auth_head, "ApiKey "), " ")
	if apiKey == "" {
		return "", fmt.Errorf("Unauthorized requested - no bearer token supplied %v", apiKey)
	}
	return apiKey, nil
}
