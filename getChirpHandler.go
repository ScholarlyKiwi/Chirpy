package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int

	if req.Method != http.MethodGet {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		chirps, err := cfg.dbq.GetChirps(req.Context())
		if err != nil {
			respBody = jsonError{Error: fmt.Sprintf("Error getting chirps: %v", err)}
			respStatus = http.StatusMethodNotAllowed
		} else {

			var jsonChirps []jsonChirp
			for _, chirp := range chirps {
				jsonChirps = append(jsonChirps, jsonChirp{
					ID:        chirp.ID,
					CreatedAt: chirp.CreatedAt,
					UpdatedAt: chirp.UpdatedAt,
					Body:      chirp.Body,
					UserID:    chirp.UserID,
				})
			}
			respBody = jsonChirps
			respStatus = http.StatusOK
		}
	}

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) getChirpByIDHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int

	respBody = ""
	if req.Method != http.MethodGet {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		chirp_id, err := uuid.Parse(req.PathValue("chirpID"))
		if err != nil {
			respStatus = http.StatusNotFound
		} else {
			chirp, err := cfg.dbq.GetChirpByID(req.Context(), chirp_id)
			if err != nil {
				respStatus = http.StatusNotFound
			} else {
				respBody = jsonChirp{
					ID:        chirp.ID,
					CreatedAt: chirp.CreatedAt,
					UpdatedAt: chirp.UpdatedAt,
					Body:      chirp.Body,
					UserID:    chirp.UserID}
				respStatus = http.StatusOK
			}
		}
	}

	jsonHtttpSend(respStatus, respBody, respWriter)
}
