package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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
	users, err := cfg.db.GetOutputUsers()
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

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
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

	tokenHeader := r.Header.Get("Authorization")
	token := tokenHeader[7:]
	user, err := cfg.authenticateToken(token)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if len(request.Email) > 0 {
		user.Email = request.Email
	}
	if len(request.Password) > 0 {
		passwordByte := []byte(request.Password)
		encryptedPW, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user.Password = string(encryptedPW)
	}
	err = cfg.db.SaveUser(user)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	safeUsers, err := cfg.db.GetOutputUsers()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, safeUsers[user.ID])
}

func (cfg *apiConfig) makeUserRed(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	var token string
	if len(tokenHeader) > 7 {
		token = tokenHeader[7:]
	}
	log.Println(token)
	if token != cfg.apiKey {
		jsonError(w, 401, "Invalid Key")
	}
	var request struct {
		Event string         `json:"event"`
		Data  map[string]int `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	if request.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "")
		return
	}

	users, _ := cfg.db.GetUsers()
	user := users[request.Data["user_id"]]
	user.IsChirpyRed = true
	cfg.db.SaveUser(user)
}
