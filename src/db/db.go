package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Connection string
}

type Connection struct {
	read  *sql.DB
	write *sql.DB
}

func (db *DB) MakeConnections() (*Connection, error) {
	read, err := sql.Open("sqlite3", db.Connection)
	if err != nil {
		return nil, err
	}
	write, err := sql.Open("sqlite3", db.Connection)
	if err != nil {
		return nil, err
	}
	log.Printf("DB Connections Made Successfully To %s", db.Connection)
	return &Connection{read, write}, nil
}
