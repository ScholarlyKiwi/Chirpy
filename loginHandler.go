package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(respWriter http.ResponseWriter, req *http.Request) {

	respBody, respStatus := cfg.login(req)

	jsonHtttpSend(respStatus, respBody, respWriter)

}

func (cfg *apiConfig) login(req *http.Request) (any, int) {
	var respBody any
	var respStatus int
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	var reqBody jsonLoginBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respBody = jsonError{Error: "Invalid JSON"}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	if len(reqBody.Email) < 3 {
		respBody = jsonError{Error: "Email Address is required"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}
	if len(reqBody.Password) < 1 {
		respBody = jsonError{Error: "Password is required"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}
	if reqBody.ExpiresInSeconds > int(time.Hour) || reqBody.ExpiresInSeconds < 1 {
		reqBody.ExpiresInSeconds = int(time.Hour)
	}

	user, err := cfg.dbq.GetUserByEmail(req.Context(), reqBody.Email)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error hashing password: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	correct, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error hashing password: %v\n", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	if !correct {
		respBody = jsonError{Error: "Password mismatch"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Token Creation error: %v.\n", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	refreshString := auth.MakeRefreshToken()
	refreshToken, err := cfg.dbq.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshString,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
	})
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Refresh Token Creation error: %v\n", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	if refreshToken.Token != refreshString {
		respBody = jsonError{Error: fmt.Sprintf("Refresh Token Creation error for token %v: %v\n", tokenString, err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	respBody = jsonUser{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefershToken: refreshString}
	respStatus = http.StatusOK

	return respBody, respStatus
}
