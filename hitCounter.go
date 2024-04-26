package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) reportHits(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("<html>	<body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", cfg.fileserverHits)
	w.Write([]byte(body))
}

func (cfg *apiConfig) hitCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Write([]byte("Hits Reset"))
}
