package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Connection   string
	Vacuum       string
	Foreign_Keys bool
}

type Connection struct {
	read  *sql.DB
	write *sql.DB
}

func (db *DB) Initalize() error {
	if _, err := os.Stat(db.Connection); errors.Is(err, os.ErrNotExist) {
		log.Println("Database file does not exist, creating...")
		_, err := os.Create(db.Connection)
		if err != nil {
			log.Fatal("Error Creating DB File")
		}
	}
	new_db, err := sql.Open("sqlite3", db.Connection)
	if err != nil {
		defer new_db.Close()
		return err
	}
	err = new_db.Ping()
	if err != nil {
		log.Fatal("DB Not Responding")
	}
	runMigrations(new_db)
	return nil
}

func runMigrations(db *sql.DB) {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("./src/db/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No Migrations To Run")
		} else {
			log.Fatal(err)
		}
	}
}

func (db *DB) connection_string() string {
	return fmt.Sprintf("%s?_auto_vacuum=%s&_foreign_keys=%t", db.Connection, db.Vacuum, db.Foreign_Keys)
}

func (db *DB) MakeConnections() (*Connection, error) {
	read, err := sql.Open("sqlite3", db.connection_string())
	if err != nil {
		return nil, err
	}
	write, err := sql.Open("sqlite3", db.connection_string())
	if err != nil {
		return nil, err
	}
	log.Printf("DB Connections Made Successfully")
	return &Connection{read, write}, nil
}
