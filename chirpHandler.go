package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/Gambor27/Chirpy/internal/database"
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
	authorID := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if len(authorID) > 0 {
		id, err := strconv.Atoi(authorID)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		authorChirps := make(map[int]database.Chirp)
		for _, chirp := range chirps {
			if chirp.AuthorID == id {
				authorChirps[chirp.ID] = chirp
			}
		}
		if len(sort) > 0 {
			sortedChirps, err := cfg.sortChirps(authorChirps, sort)
			if err != nil {
				jsonError(w, http.StatusNotFound, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, sortedChirps)
			return
		}
		respondWithJSON(w, http.StatusOK, authorChirps)
		return
	}
	if len(sort) > 0 {
		sortedChirps, err := cfg.sortChirps(chirps, sort)
		if err != nil {
			jsonError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, sortedChirps)
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

func (cfg *apiConfig) sortChirps(chirps map[int]database.Chirp, toSort string) ([]database.Chirp, error) {
	outputChirps := make([]database.Chirp, 0)
	for _, chirp := range chirps {
		outputChirps = append(outputChirps, chirp)
	}
	sortFunction, err := cfg.getSort(outputChirps, toSort)
	if err != nil {
		return make([]database.Chirp, 0), err
	}
	if len(outputChirps) > 0 {
		sort.Slice(outputChirps, sortFunction)
	}
	return outputChirps, nil
}

func (cfg *apiConfig) getSort(chirps []database.Chirp, direction string) (func(i, j int) bool, error) {
	if direction == "asc" {
		return func(i int, j int) bool {
			return chirps[i].ID < chirps[j].ID
		}, nil
	} else if direction == "desc" {
		return func(i int, j int) bool {
			return chirps[i].ID > chirps[j].ID
		}, nil
	} else {
		return nil, errors.New("sort type not known")
	}
}
