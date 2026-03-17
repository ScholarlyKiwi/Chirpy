package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
	"github.com/google/uuid"
)

type jsonCreateChirp struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(respWriter http.ResponseWriter, req *http.Request) {

	respBody, respStatus := cfg.processChirpRequest(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) processChirpRequest(req *http.Request) (respBody any, respStatus int) {
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}
	var reqBody jsonCreateChirp
	var ok bool
	var cleanedBody string
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respBody = jsonError{Error: "Invalid JSON"}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respBody = jsonError{Error: "Unable to retrieve user token."}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	uuid, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Invalid user token: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	user, err := cfg.dbq.GetUserByID(req.Context(), uuid)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Unable to retrieve user %v.", reqBody.User_id)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	cleanedBody, ok, respBody, respStatus = validateChirp(reqBody.Body)
	if ok {
		createdChrip, err := cfg.dbq.CreateChirp(req.Context(), database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: user.ID,
		},
		)
		if err != nil {
			respBody = jsonError{Error: fmt.Sprintf("Unable to create chirp: %v.", err)}
			respStatus = http.StatusBadRequest
			return respBody, respStatus
		} else {

			respBody = jsonChirp{
				ID:        createdChrip.ID,
				CreatedAt: createdChrip.CreatedAt,
				UpdatedAt: createdChrip.UpdatedAt,
				Body:      createdChrip.Body,
				UserID:    createdChrip.UserID,
			}
			respStatus = http.StatusCreated
		}
	}
	return respBody, respStatus
}
