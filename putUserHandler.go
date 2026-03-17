package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
)

func (cfg *apiConfig) putUserHandler(respWriter http.ResponseWriter, req *http.Request) {

	respBody, respStatus := cfg.putUser(req)

	jsonHtttpSend(respStatus, respBody, respWriter)

}

func (cfg *apiConfig) putUser(req *http.Request) (respBody any, respStatus int) {

	if req.Method != http.MethodPut {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
		return respBody, respStatus
	}

	bearer_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error getting token: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	uuid, err := auth.ValidateJWT(bearer_token, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Invalid user token: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	var reqBody jsonNewUser
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		respBody = jsonError{Error: "Invalid JSON"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	updatedUser, err := cfg.dbq.UpdateUser(req.Context(), database.UpdateUserParams{
		HashedPassword: hashedPassword,
		Email:          reqBody.Email,
		ID:             uuid,
	})
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error updating user: %v", err)}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	if updatedUser.Email != reqBody.Email || updatedUser.HashedPassword != hashedPassword {
		respBody = jsonError{Error: "Error updating user, value mismatch."}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}

	respBody = jsonUser{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
		Token:     bearer_token,
	}
	respStatus = http.StatusOK

	return respBody, respStatus
}
