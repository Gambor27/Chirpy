package database

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(body string, id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	currentData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		ID:       len(currentData.Chirps) + 1,
		AuthorID: id,
		Body:     body,
	}

	currentData.Chirps[chirp.ID] = chirp

	err = db.writeDB(currentData)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() (map[int]Chirp, error) {
	chirps, err := db.loadDB()
	if err != nil {
		return map[int]Chirp{}, err
	}
	return chirps.Chirps, nil
}
