package main

import (
	"fmt"
	"log"
	"net/http"
)

func readinessHandler(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	_, err := respWriter.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("Error writing readiness: %v", err)
	}
}

func (cfg *apiConfig) metricsHandler(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "text/html")
	respWriter.WriteHeader(http.StatusOK)

	body := "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n   </body>\n</html>"
	_, err := respWriter.Write([]byte(fmt.Sprintf(body, cfg.fileserverHits.Load())))
	if err != nil {
		log.Printf("Error writing readiness: %v", err)
	}
}

func (cfg *apiConfig) resetHandler(respWriter http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	cfg.dbq.DeleteUsers(req.Context())

	respWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respWriter.WriteHeader(http.StatusOK)
	_, err := respWriter.Write([]byte("File Server Hits Reset"))
	if err != nil {
		log.Printf("Error writing readiness: %v", err)
	}
}
