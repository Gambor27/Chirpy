package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) newUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		println(err)
		return
	}

	user, err := cfg.db.CreateUser(request.Email, request.Password)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) readUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.db.GetUsers()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := r.PathValue("userID")
	if len(userID) > 0 {
		userNum, chirpErr := strconv.Atoi(userID)
		if chirpErr != nil {
			jsonError(w, http.StatusBadRequest, chirpErr.Error())
			return
		}
		if users[userNum].ID == 0 {
			jsonError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithJSON(w, http.StatusOK, users[userNum])
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

/*
func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

}
*/
