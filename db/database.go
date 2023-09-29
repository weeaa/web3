package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"os"
)

type DB struct {
	db *pg.DB
}

func New() (*DB, error) {
	conn := pg.Connect(&pg.Options{
		User:     os.Getenv("PSQL_USERNAME"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Addr:     fmt.Sprintf("%s:%s", "localhost", os.Getenv("PSQL_PORT")),
		Database: os.Getenv("PSQL_DB_NAME"),
		PoolSize: 50,
	})

	if conn == nil {
		return nil, fmt.Errorf("cannot connect to postgres database")
	}

	return &DB{db: conn}, nil
}

func (pg *DB) Disconnect() {
	pg.db.Close()
}
