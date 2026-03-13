package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(respWriter http.ResponseWriter, req *http.Request) {

	respBody, respStatus := cfg.login(respWriter, req)

	jsonHtttpSend(respStatus, respBody, respWriter)

}

func (cfg *apiConfig) login(respWriter http.ResponseWriter, req *http.Request) (any, int) {
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

	user, err := cfg.dbq.GetUserByEmail(req.Context(), reqBody.Email)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error hashing password: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	correct, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Error hashing password: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	if !correct {
		respBody = jsonError{Error: "Password mismatch"}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus
	}

	respBody = jsonUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdateAt,
		Email:     user.Email}
	respStatus = http.StatusOK

	return respBody, respStatus
}
