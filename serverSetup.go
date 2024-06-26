package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/Gambor27/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	secret         string
	apiKey         string
}

func serverSetup() error {
	chirpDB, err := database.NewDB("./db")
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             chirpDB,
		secret:         os.Getenv("JWT_SECRET"),
		apiKey:         os.Getenv("API_KEY"),
	}
	mux.Handle("/app/*", apiCfg.hitCounter(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/reset/", apiCfg.resetHits)
	mux.HandleFunc("/", directory)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.readChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.readChirps)
	mux.HandleFunc("GET /api/healthz", health)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.reportHits)
	mux.HandleFunc("GET /api/users", apiCfg.readUsers)
	mux.HandleFunc("GET /api/users/{userID}", apiCfg.readUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.newChirp)
	mux.HandleFunc("POST /api/login", apiCfg.authLogin)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.makeUserRed)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeToken)
	mux.HandleFunc("POST /api/users", apiCfg.newUser)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}
	err = server.ListenAndServe()
	if err != nil {
		return errors.New("server failed to start")
	}
	return nil
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func directory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><a href=\"./app\">Homepage</a>\n<a href=\"./api/healthz\">Health</a>\n<a href=\"./api/reset\">Reset</a>\n<a href=\"./admin/metrics/\">Metrics</a></body></html>"))
}
