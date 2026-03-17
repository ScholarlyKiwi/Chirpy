package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ScholarlyKiwi/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbq            *database.Queries
	tokenSecret    string
}

const filepathRoot = "."
const port = "8080"

func main() {

	var apiCfg apiConfig

	err := getConfig(&apiCfg)
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
		return
	}

	serveMux := http.NewServeMux()

	apiCfg.assignHandlers(serveMux)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func (apiCfg *apiConfig) assignHandlers(serveMux *http.ServeMux) {
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.emailHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.chirpHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirpHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	serveMux.HandleFunc("PUT /api/users", apiCfg.putUserHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpHandler)
}

func getConfig(cfg *apiConfig) error {
	err := godotenv.Load()

	log.Default()
	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok || dbURL == "" {
		return fmt.Errorf("Missing env DB_URL")
	}
	severToken, ok := os.LookupEnv("SECRET")
	if !ok || severToken == "" {
		return fmt.Errorf("Missing env SECRET")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return fmt.Errorf("Error accessing database: %v", err)
	}

	cfg.dbq = database.New(db)
	cfg.tokenSecret = severToken
	return nil
}
