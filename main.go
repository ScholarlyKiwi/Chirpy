package main

import (
	"database/sql"
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
}

const filepathRoot = "."
const port = "8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	log.Default()
	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok || dbURL == "" {
		log.Fatalf("Missing DB_URL")
		return
	}
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Error accessing database: %v", err)
		return
	}

	var apiCfg apiConfig
	apiCfg.dbq = database.New(db)

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
}
