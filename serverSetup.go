package main

import (
	"errors"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	chirpID        int
}

func serverSetup() error {

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: 0,
		chirpID:        0,
	}
	mux.Handle("/app/*", apiCfg.hitCounter(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/", directory)
	mux.HandleFunc("GET /api/healthz", health)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.reportHits)
	mux.HandleFunc("/api/reset/", apiCfg.resetHits)
	mux.HandleFunc("POST /api/chirps", apiCfg.newChirp)
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}
	err := server.ListenAndServe()
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
