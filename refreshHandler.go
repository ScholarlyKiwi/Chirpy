package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(respWriter http.ResponseWriter, req *http.Request) {
	respBody, respStatus := cfg.refresh(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) refresh(req *http.Request) (any, int) {
	var respBody any
	var respStatus int
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	bearer_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error getting token: %v", err)}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}
	refreshUser, err := cfg.dbq.GetUserByRefreshToken(req.Context(), bearer_token)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Refresh token not fond: %v", err)}
		respStatus = 401
		return respBody, respStatus
	}
	if time.Now().After(refreshUser.ExpiresAt) {
		respBody = jsonError{Error: fmt.Sprintf("Refresh token has expired at: %v", refreshUser.ExpiresAt)}
		respStatus = 401
		return respBody, respStatus
	}
	if refreshUser.RevokedAt.Valid {
		respBody = jsonError{Error: fmt.Sprintf("Refresh token was revoked at: %v", refreshUser.RevokedAt.Time)}
		respStatus = 401
		return respBody, respStatus
	}

	newtokenString, err := auth.MakeJWT(refreshUser.ID, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Token Creation error: %v.\n", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	respBody = jsonToken{
		Token: newtokenString,
	}
	respStatus = http.StatusOK

	return respBody, respStatus
}
