package main

import (
	"errors"
	"net/http"

	"github.com/Gambor27/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
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
	}
	mux.Handle("/app/*", apiCfg.hitCounter(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	//mux.HandleFunc("/", directory)
	mux.HandleFunc("GET /api/healthz", health)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.reportHits)
	mux.HandleFunc("/api/reset/", apiCfg.resetHits)
	mux.HandleFunc("POST /api/chirps", apiCfg.newChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.readChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.readChirps)
	mux.HandleFunc("POST /api/users", apiCfg.newUser)
	mux.HandleFunc("GET /api/users", apiCfg.readUsers)
	mux.HandleFunc("GET /api/users/{userID}", apiCfg.readUsers)
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

//func directory(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("<html><body><a href=\"./app\">Homepage</a>\n<a href=\"./api/healthz\">Health</a>\n<a href=\"./api/reset\">Reset</a>\n<a href=\"./admin/metrics/\">Metrics</a></body></html>"))
//}
