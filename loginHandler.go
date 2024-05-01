package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Gambor27/Chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) authLogin(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
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
	var loginResponse struct {
		User  database.OutputUser
		Token string `json:"token"`
	}
	expiration := 86400
	if request.Expires_in_seconds > 0 && request.Expires_in_seconds < expiration {
		expiration = request.Expires_in_seconds
	}

	loginResponse.User = user
	loginResponse.Token, err = cfg.createToken(user, expiration)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, loginResponse)
}

func (cfg *apiConfig) createToken(user database.OutputUser, expiration int) (string, error) {

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiration))),
		Subject:   strconv.Itoa(user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(cfg.secret)
	output, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return output, nil
}
