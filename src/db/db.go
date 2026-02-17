// Package db provides functionality for interacting with the application's database.
package db

import (
	"BloTils/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3" // Ensure SQLite driver is registered for database/sql
)

var new_db *sql.DB

// closeDB safely closes the provided database connection.
// It defers the closing operation and logs any error encountered during the closure.
// This function helps ensure that database resources are properly released.
func CloseDB() {
	defer func() {
		if err := new_db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
		log.Println("Database connection closed")
	}()
}

// closeFile safely closes the provided source.Driver, logging any errors encountered during the close operation.
// It uses a deferred function to ensure that the Close method is called when closeFile returns.
func closeFile(fSrc source.Driver) {
	defer func() {
		if err := fSrc.Close(); err != nil {
			log.Printf("Error Closing Migration Files: %s", err)
		}
	}()
}

// InitDB initializes the database connection using the provided DBConfig.
// It checks if the database file exists, and if not, logs that it will be created.
// The function attempts to open a SQLite database connection and pings it to ensure responsiveness.
// If the connection is successful, it runs database migrations using the specified migration files.
// On success, the established connection is stored in the DBConfig; otherwise, errors are logged and returned.
//
// Parameters:
//
//	db - Pointer to a models.DBConfig struct containing database configuration.
//
// Returns:
//
//	error - An error if the connection or migrations fail, otherwise nil.
func InitDB(db *models.DBConfig) error {
	if _, err := os.Stat(db.DBLocation); errors.Is(err, os.ErrNotExist) {
		log.Println("Database file does not exist, will be created by SQLite on connect...")
	}
	new_db, err := sql.Open("sqlite3", connectionString(db))
	if err != nil {
		return err
	}
	err = new_db.Ping()
	if err != nil {
		log.Println("DB Not Responding")
		if cerr := new_db.Close(); cerr != nil {
			log.Println("Error closing database connection:", cerr)
		}
		return err
	}
	err = runMigrations(new_db, db.MigrationFiles)
	if err != nil {
		log.Println("Error Running Migrations")
		CloseDB()
		db.Connection = nil
		return err
	}
	db.Connection = new_db
	return nil
}

// runMigrations applies database schema migrations to the provided SQLite database connection.
// It takes a *sql.DB instance and a path to the migration files as arguments.
// The function initializes a migration instance using the provided migration files and database connection,
// then attempts to apply all up migrations. If there are no migrations to run, it logs this information.
// Any errors encountered during the process are logged and returned.
// Resources are properly closed after use.
func runMigrations(db *sql.DB, migrationFiles string) error {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Printf("Error Connecting With SQLite Instance: %s", err)
		return err
	}
	fSrc, err := (&file.File{}).Open(migrationFiles)
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

// connectionString constructs a SQLite connection string using the provided DBConfig.
// It formats the connection string to include auto vacuum and foreign key options.
//
// Parameters:
//
//	db - a pointer to a models.DBConfig struct containing database configuration.
//
// Returns:
//
//	A string representing the formatted SQLite connection string.
func connectionString(db *models.DBConfig) string {
	return fmt.Sprintf("%s?_auto_vacuum=%s&_foreign_keys=%t", db.DBLocation, db.Vacuum, db.ForeignKeys)
}
