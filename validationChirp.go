package main

import (
	"net/http"
	"slices"
	"strings"
)

const maxChirpLength = 140

func validateChirp(chirpBody string) (cleanedBody string, ok bool, respErr jsonError, statusCode int) {

	if len(chirpBody) > maxChirpLength {
		respErr = jsonError{Error: "Chirp is too long"}
		statusCode = http.StatusBadRequest
		ok = false
	} else {
		cleanedBody = cleanBody(chirpBody)
		statusCode = http.StatusAccepted
		ok = true
	}
	return cleanedBody, ok, respErr, statusCode
}

func cleanBody(unclean string) string {
	unclean_words := []string{"kerfuffle", "sharbert", "fornax"}
	clean_words := []string{}

	for _, word := range strings.Split(unclean, " ") {
		if slices.Contains(unclean_words, strings.ToLower(word)) {
			word = "****"
		}
		clean_words = append(clean_words, word)
	}

	return strings.Join(clean_words, " ")
}
