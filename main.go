package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"log"
	"net/http"
	"sync/atomic"
	"os"

	"github.com/rangaroo/chirpy-http-server/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could't open the database: %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	server := &http.Server{
		Addr:     ":" + port,
		Handler:  mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
