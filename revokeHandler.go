package main

import (
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(respWriter http.ResponseWriter, req *http.Request) {
	respBody, respStatus := cfg.revoke(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) revoke(req *http.Request) (respBody any, respStatus int) {
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	bearer_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error getting token: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	refreshToken, err := cfg.dbq.GetRefreshTokenByToken(req.Context(), bearer_token)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error retrieving refresh token: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	if refreshToken.RevokedAt.Valid {
		respBody = jsonError{Error: fmt.Sprintf("Refresh Token already Revoked: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	revokedToken, err := cfg.dbq.RevokeRefreshToken(req.Context(), refreshToken.Token)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error revoking token: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	if !revokedToken.RevokedAt.Valid {
		respBody = jsonError{Error: "Error revoking token - token not revoked"}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	return respBody, http.StatusNoContent
}
