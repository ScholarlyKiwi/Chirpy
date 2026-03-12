package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/google/uuid"
)

type jsonEmailBody struct {
	Email string `json:"email"`
}

type jsonUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) emailHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		var reqBody jsonEmailBody
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			respBody = jsonError{Error: "Invalid JSON"}
			respStatus = http.StatusBadRequest
		} else if len(reqBody.Email) < 3 {
			respBody = jsonError{Error: "Email Address is required"}
			respStatus = http.StatusBadRequest
		} else {
			user, err := cfg.dbq.CreateUser(req.Context(), reqBody.Email)

			if err != nil {
				respBody = jsonError{Error: fmt.Sprintf("Unable to create User: %v", err)}
				respStatus = http.StatusBadRequest
			}
			respBody = jsonUser{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdateAt,
				Email:     user.Email}
			respStatus = http.StatusCreated
		}
	}

	jsonHtttpSend(respStatus, respBody, respWriter)

}
