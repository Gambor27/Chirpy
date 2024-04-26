package main

import (
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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
