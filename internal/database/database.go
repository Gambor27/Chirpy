package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func NewDB(path string) (*DB, error) {
	mux := &sync.RWMutex{}
	db := DB{
		path: path,
		mux:  mux,
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, fs.ErrNotExist) {
		file, createErr := os.Create(db.path)
		if createErr != nil {
			//Error occurred while attempting to create file.  Return error and exit method.
			return createErr
		} else {
			defer file.Close()
			//File did not exist.  Successfully created.  Exit method with no error.
			return nil
		}

	} else if err != nil {
		//Reading file returned error other than not exists.  Return error and exit method.
		return err
	}
	//File already exists.  Exit method with no error.
	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	currentData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		ID:   len(currentData.Chirps) + 1,
		Body: body,
	}

	currentData.Chirps[chirp.ID] = chirp

	err = db.writeDB(currentData)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	currentData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:    len(currentData.Chirps) + 1,
		Email: email,
	}

	currentData.Users[user.ID] = user

	err = db.writeDB(currentData)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetChirps() (map[int]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	chirps, err := db.loadDB()
	if err != nil {
		return map[int]Chirp{}, err
	}
	return chirps.Chirps, nil
}

func (db *DB) GetUsers() (map[int]User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dat, err := db.loadDB()
	if err != nil {
		return map[int]User{}, err
	}
	return dat.Users, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	contents, err := os.ReadFile(db.path)
	if err != nil {
		log.Println(err)
		return DBStructure{}, err
	}
	dat := make(map[int]Chirp)
	output := DBStructure{
		Chirps: dat,
	}

	if len(contents) > 0 {
		err = json.Unmarshal(contents, &output)
		if err != nil {
			log.Println(err)
			return DBStructure{}, err
		}
	}
	return output, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	fileContents, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, fileContents, 0666)
	if err != nil {
		return err
	}
	return nil
}
