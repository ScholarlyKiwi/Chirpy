package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webhooksHandler(respWriter http.ResponseWriter, req *http.Request) {
	respBody, respStatus := cfg.processWebhooks(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

var eventTypes = []string{
	"user.upgraded",
}

func (cfg *apiConfig) processWebhooks(req *http.Request) (respBody any, respStatus int) {

	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	var hook jsonWebhook
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&hook)
	if err != nil {
		respBody = jsonError{Error: "Invalid JSON"}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	if !slices.Contains(eventTypes, hook.Event) {
		respBody = jsonError{Error: "Inavlid webhook event."}
		respStatus = http.StatusNoContent
		return respBody, respStatus
	}

	user_id, err := uuid.Parse(hook.Data.UserID)
	if err != nil {
		respBody = jsonError{Error: "Invalid User ID"}
		respStatus = http.StatusNotFound
		return respBody, respStatus
	}
	user, err := cfg.dbq.GetUserByID(req.Context(), user_id)
	if err != nil {
		respBody = jsonError{Error: "User ID not found"}
		respStatus = http.StatusNotFound
		return respBody, respStatus
	}

	switch hook.Event {
	case "user.upgraded":
		respBody, respStatus = cfg.processUserUpgrade(req, user.ID)
	default:
		respBody = jsonError{Error: "Webhook event not implemented."}
		respStatus = http.StatusBadRequest
	}

	return respBody, respStatus
}

func (cfg *apiConfig) processUserUpgrade(req *http.Request, userID uuid.UUID) (respBody any, respStatus int) {

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respBody = jsonError{Error: "Missing API Key"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}
	if apiKey != cfg.polkaKey {
		respBody = jsonError{Error: "Invalid API Key supplied"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	_, err = cfg.dbq.UpdateUserIsChirpy(req.Context(), database.UpdateUserIsChirpyParams{
		IsChirpyRed: true,
		ID:          userID},
	)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error updating IsChirpyRed: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	respBody = ""
	respStatus = http.StatusNoContent

	return respBody, respStatus
}
