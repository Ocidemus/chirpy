package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

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

func(cfg *apiConfig) reset(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK) // 200
	cfg.fileserverHits.Store(0)
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w , r)
	})
}

func main(){
	const filepathRoot = "."
	const port = "8080"
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", endpoint)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot)))))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

