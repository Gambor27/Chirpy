package database

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int                  `json:"id"`
	Email        string               `json:"email"`
	Password     string               `json:"password"`
	RefreshToken map[string]time.Time `json:"refreshToken"`
}

type OutputUser struct {
	ID           int                  `json:"id"`
	Email        string               `json:"email"`
	RefreshToken map[string]time.Time `json:"refreshToken,omitempty"`
}

func (db *DB) CreateUser(email string, password string) (OutputUser, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	currentData, err := db.loadDB()
	if err != nil {
		return OutputUser{}, err
	}
	passwordByte := []byte(password)
	encryptedPW, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return OutputUser{}, err
	}
	passwordString := string(encryptedPW)

	for _, user := range currentData.Users {
		if user.Email == email {
			return OutputUser{}, errors.New("user already exists")
		}
	}

	user := User{
		ID:       len(currentData.Users) + 1,
		Email:    email,
		Password: passwordString,
	}

	currentData.Users[user.ID] = user

	err = db.writeDB(currentData)
	if err != nil {
		return OutputUser{}, err
	}

	return db.getOutputUser(user), nil
}

func (db *DB) GetUsers() (map[int]User, error) {
	dat, err := db.loadDB()
	if err != nil {
		return map[int]User{}, err
	}
	return dat.Users, nil
}

func (db *DB) GetOutputUsers() (map[int]OutputUser, error) {
	dat, err := db.loadDB()
	if err != nil {
		return map[int]OutputUser{}, err
	}
	outputMap := make(map[int]OutputUser)
	for _, user := range dat.Users {
		safeuser := db.getOutputUser(user)
		outputMap[safeuser.ID] = safeuser
	}
	return outputMap, nil
}

func (db *DB) ValidateUser(username, password string) (OutputUser, error) {
	users, err := db.GetUsers()
	if err != nil {
		return OutputUser{}, err
	}
	requestPasswordByte := []byte(password)
	for _, user := range users {
		if user.Email == username {
			userPasswordByte := []byte(user.Password)
			err := bcrypt.CompareHashAndPassword(userPasswordByte, requestPasswordByte)
			if err != nil {
				return OutputUser{}, err
			}
			return db.getOutputUser(user), nil
		}
	}
	return OutputUser{}, errors.New("failed to process login attempt")
}

func (db *DB) getOutputUser(user User) OutputUser {
	return OutputUser{
		ID:           user.ID,
		Email:        user.Email,
		RefreshToken: user.RefreshToken,
	}
}

func (db *DB) SaveUser(user User) error {
	dat, err := db.loadDB()
	if err != nil {
		return err
	}
	dat.Users[user.ID] = user
	db.mux.Lock()
	defer db.mux.Unlock()
	db.writeDB(dat)
	return nil
}

func (db *DB) UpdateRefreshToken(token string, userid int) error {
	users, err := db.GetUsers()
	if err != nil {
		return err
	}
	tokenMap := make(map[string]time.Time)
	tokenMap[token] = time.Now().Add(time.Hour * 1440)
	user := users[userid]
	user.RefreshToken = tokenMap
	err = db.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}
