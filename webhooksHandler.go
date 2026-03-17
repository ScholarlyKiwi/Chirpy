package main

import (
	"net/http"
)

func (cfg *apiConfig) webhooksHandler(respWriter http.ResponseWriter, req *http.Request) {
	respBody, respStatus := cfg.processWebhooks(req)

	jsonHtttpSend(respStatus, respBody, respWriter)
}

func (cfg *apiConfig) processWebhooks(req *http.Request) (respBody any, respStatus int) {
	return respBody, respStatus
}
