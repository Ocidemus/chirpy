package main

import (
	
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	

	"database/sql"

	"github.com/Ocidemus/chirpy/internal/database"
	
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func endpoint(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte("OK"))

}

func(cfg *apiConfig) metrics(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK) // 200
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,hits)))
}


type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w , r)
	})
}

func main(){
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")	
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
    	log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == ""{
		log.Fatal("Secret must be set")
	}
	dbconn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	


	const filepathRoot = "."
	const port = "8080"
	cfg := &apiConfig{
	platform: platform,
	secret: secret,
}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", endpoint)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("GET /api/chirps",cfg.reqchirp)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/users", cfg.create_user)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}",cfg.handlechirp )
	mux.HandleFunc("POST /api/login",cfg.handle_login)
	mux.HandleFunc("POST /api/refresh",cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke",cfg.handleRevoke)

	// mux.HandleFunc(("POST /api/chirp"))
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot)))))
	dbQueries := database.New(dbconn)
	cfg.db = dbQueries
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

