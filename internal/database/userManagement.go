package database

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	currentData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:    len(currentData.Users) + 1,
		Email: email,
	}

	currentData.Users[user.ID] = user

	err = db.writeDB(currentData)
	if err != nil {
		return User{}, err
	}

	return user, nil
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
