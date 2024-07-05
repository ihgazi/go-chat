package db

import (
	"database/sql"
	_ "github.com/lib/pq"

    "github.com/ihgazi/go-chat/config"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
    PostgresUrl := config.LoadConfig().PostgresUrl

	db, err := sql.Open("postgres", PostgresUrl)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
