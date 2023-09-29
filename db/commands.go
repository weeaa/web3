package db

type DatabaseService interface {
	InsertUser(*User) error
	GetUser(string) (*User, error)
	UpdateUser(string) error
	RemoveUser(string) error
}

func (pg *DB) InsertUser(u *User) error {
	_, err := pg.db.Model(u).Insert()
	return err
}

func (pg *DB) GetUser(baseAddress string) (*User, error) {
	return nil, nil
}

func (pg *DB) UpdateUser(baseAddress string) error {
	return nil
}

func (pg *DB) RemoveUser(baseAddress string) error {
	_, err := pg.db.Model(User{BaseAddress: baseAddress}).WherePK().Delete()
	return err
}
