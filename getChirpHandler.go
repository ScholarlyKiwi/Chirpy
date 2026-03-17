package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/ScholarlyKiwi/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpHandler(respWriter http.ResponseWriter, req *http.Request) {
	var respBody any
	var respStatus int

	if req.Method != http.MethodGet {
		respBody = jsonError{Error: "Invalid request method"}
		respStatus = http.StatusMethodNotAllowed
	} else {
		author_id := req.URL.Query().Get("author_id")
		if author_id != "" {
			respBody, respStatus = cfg.getChirpsByAuthorID(req, author_id)
		} else {
			chirps, err := cfg.dbq.GetChirps(req.Context())
			if err != nil {
				respBody = jsonError{Error: fmt.Sprintf("Error getting chirps: %v", err)}
				respStatus = http.StatusMethodNotAllowed
			} else {
				chirps = sortChirps(chirps, req)
				respBody = convertChirps(chirps)
				respStatus = http.StatusOK
			}
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

func (cfg *apiConfig) getChirpsByAuthorID(req *http.Request, author_id string) (respBody any, respStatus int) {

	author_uuid, err := uuid.Parse(author_id)
	if err != nil {
		respBody = jsonError{Error: "Invalid Author ID"}
		respStatus = http.StatusBadRequest
		return respBody, respStatus
	}
	author, err := cfg.dbq.GetUserByID(req.Context(), author_uuid)
	if err != nil {
		respBody = jsonError{Error: "Author ID not found"}
		respStatus = http.StatusNotFound
		return respBody, respStatus
	}
	chirps, err := cfg.dbq.GetChirpsByUserID(req.Context(), author.ID)
	if err != nil {
		respBody = jsonError{Error: "No Chirps found"}
		respStatus = http.StatusNotFound
		return respBody, respStatus
	}

	chirps = sortChirps(chirps, req)
	respBody = convertChirps(chirps)
	respStatus = http.StatusOK

	return respBody, respStatus
}

func convertChirps(dbChirps []database.Chirp) []jsonChirp {

	var jsonChirps []jsonChirp
	for _, chirp := range dbChirps {
		jsonChirps = append(jsonChirps, jsonChirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	return jsonChirps
}

func sortChirps(dbChirps []database.Chirp, req *http.Request) []database.Chirp {
	sortMethod := req.URL.Query().Get("sort")

	sort.Slice(dbChirps,
		func(i, j int) bool {
			if sortMethod == "desc" {
				return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
			} else {
				return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt)
			}
		})
	return dbChirps
}
