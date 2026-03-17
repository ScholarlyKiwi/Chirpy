package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
)

func (cfg *apiConfig) emailHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		var reqBody jsonLoginBody
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			respBody = jsonError{Error: "Invalid JSON"}
			respStatus = http.StatusBadRequest
		} else if len(reqBody.Email) < 3 {
			respBody = jsonError{Error: "Email Address is required"}
			respStatus = http.StatusBadRequest
		} else if len(reqBody.Password) < 1 {
			respBody = jsonError{Error: "Password is required"}
			respStatus = http.StatusBadRequest
		} else {
			hashed_password, err := auth.HashPassword(reqBody.Password)
			if err != nil {
				respBody = jsonError{Error: fmt.Sprintf("Error hashing password: %v", err)}
				respStatus = http.StatusBadRequest
			} else {

				user, err := cfg.dbq.CreateUser(req.Context(), database.CreateUserParams{
					Email:          reqBody.Email,
					HashedPassword: hashed_password,
				})
				if err != nil {
					respBody = jsonError{Error: fmt.Sprintf("Unable to create User: %v", err)}
					respStatus = http.StatusBadRequest
				} else {
					respBody = jsonUser{
						ID:         user.ID,
						CreatedAt:  user.CreatedAt,
						UpdatedAt:  user.UpdatedAt,
						Email:      user.Email,
						IsChirpRed: user.IsChirpyRed}
					respStatus = http.StatusCreated
				}
			}
		}
	}

	jsonHtttpSend(respStatus, respBody, respWriter)

}
