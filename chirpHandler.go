package main

import (
	"encoding/json"
	"net/http"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) newChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	chirp.Body = profanityFilter(chirp.Body)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't decode Chirp")
		return
	} else if len(chirp.Body) > 140 {
		jsonError(w, http.StatusInternalServerError, "Chirp to long")
		return
	}

	cfg.chirpID++
	chirp.ID = cfg.chirpID
	respondWithJSON(w, http.StatusCreated, chirp)
}
