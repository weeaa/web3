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

func (pg *DB) GetUser(baseAddress string) (User, error) {
	var user User
	if err := pg.db.Model(&user).Where("base_address = ?", baseAddress).Select(); err != nil {
		return User{}, err
	}
	return user, nil
}

func (pg *DB) UpdateUser(baseAddress string) error {

	return nil
}

func (pg *DB) RemoveUser(baseAddress string) error {
	_, err := pg.db.Model(User{BaseAddress: baseAddress}).WherePK().Delete()
	return err
}

func (pg *DB) GetAllUsers() ([]User, error) {
	var users []User
	err := pg.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

/*
func (pg *DB) WipeTable() error {
	pg.db.Model().Delete()
	return nil
}
*/
