package main

import (
	"encoding/json"
	"net/http"
)

func jsonHtttpSend(respStatus int, respBody any, respWriter http.ResponseWriter) {
	var data []byte
	var err error
	data, err = json.Marshal(respBody)

	if err != nil {
		respBody = jsonError{Error: "Error Writing JSON"}
		respStatus = http.StatusInternalServerError
	}

	respWriter.Header().Set("Contetnt-Type", "application/json")
	respWriter.WriteHeader(respStatus)
	respWriter.Write(data)

}
