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
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	} else if len(request.Body) > 140 {
		jsonError(w, http.StatusBadRequest, "Chirp to long")
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	token := tokenHeader[7:]
	user, err := cfg.authenticateToken(token)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(request.Body, user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) readChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpID := r.PathValue("chirpID")
	if len(chirpID) > 0 {
		chirpNum, chirpErr := strconv.Atoi(chirpID)
		if chirpErr != nil {
			jsonError(w, http.StatusBadRequest, chirpErr.Error())
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	token := tokenHeader[7:]
	user, err := cfg.authenticateToken(token)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpNum, chirpErr := strconv.Atoi(chirpID)
	if chirpErr != nil {
		jsonError(w, http.StatusBadRequest, chirpErr.Error())
		return
	}
	chirp := chirps[chirpNum]

	if chirp.ID == 0 {
		jsonError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	if chirp.AuthorID != user.ID {
		jsonError(w, http.StatusForbidden, "Access denied")
		return
	}

	delete(chirps, chirpNum)
	err = cfg.db.SaveChirps(chirps)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "Chirp Deleted")
}
