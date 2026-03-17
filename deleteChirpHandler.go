package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(respWriter http.ResponseWriter, req *http.Request) {
	respBody, respStatus := cfg.deleteChirp(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) deleteChirp(req *http.Request) (respBody any, respStatus int) {
	if req.Method != http.MethodDelete {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respBody = jsonError{Error: "Unable to retrieve user token."}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	user_uuid, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Invalid user token: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	chirp_id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respBody = jsonError{Error: "Error getting chirp ID."}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	chirpDelete, err := cfg.dbq.GetChirpByID(req.Context(), chirp_id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			respBody = jsonError{Error: "Chirp not found"}
			respStatus = http.StatusNotFound
		} else {
			respBody = jsonError{Error: fmt.Sprintf("Error finding chirp %v", err)}
			respStatus = http.StatusUnauthorized
		}
		return respBody, respStatus
	}
	if chirpDelete.ID != chirp_id {
		respBody = jsonError{Error: "ChirpDelete: ID mismatch loading Chirp"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}
	if chirpDelete.UserID != user_uuid {
		respBody = jsonError{Error: "User is not authorised to delete this Chirp"}
		respStatus = http.StatusForbidden
		return respBody, respStatus
	}

	err = cfg.dbq.DeleteChirpByID(req.Context(), chirp_id)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error deleting chirp: %v", err)}
		respStatus = http.StatusNoContent
		return respBody, respStatus
	}
	respStatus = http.StatusNoContent

	return respBody, respStatus
}
