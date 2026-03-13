package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/database"
	"github.com/google/uuid"
)

type jsonCreateChirp struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int
	if req.Method != http.MethodPost {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		var reqBody jsonCreateChirp
		var ok bool
		var cleanedBody string
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			respBody = jsonError{Error: "Invalid JSON"}
			respStatus = http.StatusBadRequest
		}

		user, err := cfg.dbq.GetUserByID(req.Context(), reqBody.User_id)
		if err != nil {
			respBody = jsonError{Error: "Unable to retrieve user."}
			respStatus = http.StatusBadRequest
		} else {
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
		}
	}

	jsonHtttpSend(respStatus, respBody, respWriter)
}
