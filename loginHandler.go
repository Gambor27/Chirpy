package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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
		ID           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	expiration := 3600

	if request.Expires_in_seconds > 0 && request.Expires_in_seconds < expiration {
		expiration = request.Expires_in_seconds
	}

	loginResponse.Token, err = cfg.createToken(user, expiration)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	loginResponse.ID = user.ID
	loginResponse.Email = user.Email
	loginResponse.IsChirpyRed = user.IsChirpyRed
	loginResponse.RefreshToken = cfg.createRefresh()
	err = cfg.db.UpdateRefreshToken(loginResponse.RefreshToken, loginResponse.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, loginResponse)
}

func (cfg *apiConfig) authenticateToken(tokenString string) (database.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		return database.User{}, err
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return database.User{}, err
	}
	users, err := cfg.db.GetUsers()
	if err != nil {
		return database.User{}, err
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		return database.User{}, err
	}
	output, ok := users[intID]
	if !ok {
		return database.User{}, errors.New("userid no longer exists")
	}
	return output, nil
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

func (cfg *apiConfig) createRefresh() string {
	seed := make([]byte, 32)
	rand.Read(seed)
	return hex.EncodeToString(seed)
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	token := tokenHeader[7:]
	users, err := cfg.db.GetOutputUsers()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	var reply struct {
		Token string `json:"token"`
	}
	for _, user := range users {
		expiration, ok := user.RefreshToken[token]
		if ok && expiration.After(time.Now()) {
			reply.Token, err = cfg.createToken(user, 3600)
			if err != nil {
				respondWithJSON(w, http.StatusInternalServerError, err)
				return
			}
			respondWithJSON(w, http.StatusOK, reply)
			return
		}
	}
	jsonError(w, http.StatusUnauthorized, "Token Not Found")
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	token := tokenHeader[7:]
	users, err := cfg.db.GetUsers()
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, user := range users {
		expiration, ok := user.RefreshToken[token]
		if ok && expiration.After(time.Now()) {
			user.RefreshToken = nil
			cfg.db.SaveUser(user)
			respondWithJSON(w, http.StatusOK, "")
		}
	}
}
