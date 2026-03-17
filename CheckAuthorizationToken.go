package main

import (
	"fmt"
	"net/http"

	"github.com/ScholarlyKiwi/Chirpy/internal/auth"
	"github.com/ScholarlyKiwi/Chirpy/internal/database"
)

func (cfg *apiConfig) CheckAuthorizationToken(req *http.Request) (respBody any, respStatus int, userRecord database.User, err error) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Println(err)
		respBody = jsonError{Error: "Unable to retrieve user token."}
		respStatus = http.StatusBadRequest
		return respBody, respStatus, userRecord, err
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("Invalid user token: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus, userRecord, err
	}

	userRecord, err = cfg.dbq.GetUserByID(req.Context(), userID)
	if err != nil {
		respBody = jsonError{Error: fmt.Sprintf("User not found: %v", err)}
		respStatus = http.StatusUnauthorized
		return respBody, respStatus, userRecord, err
	}

	return respBody, respStatus, userRecord, err
}
