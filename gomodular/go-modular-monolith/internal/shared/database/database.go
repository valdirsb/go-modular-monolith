package database

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Connection: db}, nil
}

func (d *Database) Close() error {
	return d.Connection.Close()
}
