package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) newChirp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	request.Body = profanityFilter(request.Body)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Couldn't decode Chirp")
		return
	} else if len(request.Body) > 140 {
		jsonError(w, http.StatusBadRequest, "Chirp to long")
		return
	}

	chirp, err := cfg.db.CreateChirp(request.Body)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Error creating Chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) readChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "Failed to load Chirps")
		return
	}

	chirpID := r.PathValue("chirpID")
	if len(chirpID) > 0 {
		chirpNum, chirpErr := strconv.Atoi(chirpID)
		if chirpErr != nil {
			jsonError(w, http.StatusBadRequest, "Error retreiving ChirpID")
			return
		}
		if chirps[chirpNum].ID == 0 {
			jsonError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithJSON(w, http.StatusOK, chirps[chirpNum])
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
