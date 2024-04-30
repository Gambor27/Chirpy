package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OutputUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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
		ID:    user.ID,
		Email: user.Email,
	}
}
