package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) authLogin(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.ValidateUser(request.Email, request.Password)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		log.Println(err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
