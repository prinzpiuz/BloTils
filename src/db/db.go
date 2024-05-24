// Package db provides functionality for interacting with the application's database.
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	DBLocation  string
	Vacuum      string
	ForeignKeys bool
	Connection  *sql.DB
}

func closeDB(new_db *sql.DB) {
	defer func() {
		if err := new_db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}()
}

func closeFile(fSrc source.Driver) {
	defer func() {
		if err := fSrc.Close(); err != nil {
			log.Printf("Error Closing Migration Files: %s", err)
		}
	}()
}

func (db *DB) Initialize() error {
	if _, err := os.Stat(db.DBLocation); errors.Is(err, os.ErrNotExist) {
		log.Println("Database file does not exist, creating...")
		_, err := os.Create(db.DBLocation)
		if err != nil {
			return err
		}
	}
	new_db, err := sql.Open("sqlite3", db.connection_string())
	if err != nil {
		closeDB(new_db)
		return err
	}
	err = new_db.Ping()
	if err != nil {
		log.Println("DB Not Responding")
		closeDB(new_db)
		return err
	}
	err = runMigrations(new_db)
	if err != nil {
		log.Println("Error Running Migrations")
		closeDB(new_db)
		return err
	}
	db.Connection = new_db
	return nil
}

func runMigrations(db *sql.DB) error {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Printf("Error Connecting With SQLite Instance: %s", err)
		return err
	}

	fSrc, err := (&file.File{}).Open("./src/db/migrations")
	if err != nil {
		closeFile(fSrc)
		log.Printf("Error Getting Migration Files: %s", err)
		return err
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		closeFile(fSrc)
		log.Printf("Error Creating Migration Instance: %s", err)
		return err
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No Migrations To Run")
		} else {
			closeFile(fSrc)
			log.Printf("Error While Running UP Migrations: %s", err)
			return err
		}
	}
	closeFile(fSrc)
	return nil
}

func (db *DB) connection_string() string {
	return fmt.Sprintf("%s?_auto_vacuum=%s&_foreign_keys=%t", db.DBLocation, db.Vacuum, db.ForeignKeys)
}
