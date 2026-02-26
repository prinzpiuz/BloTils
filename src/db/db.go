// Package db provides functionality for interacting with the application's database.
package db

import (
	"BloTils/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbConnection *sql.DB
	dbMutex      sync.RWMutex
	closeOnce    sync.Once
)

// InitDB initializes the database connection using the provided DBConfig.
// It creates the database file if it doesn't exist, runs migrations,
// and stores the connection for later use.
func InitDB(config *models.DBConfig) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Check if already initialized
	if dbConnection != nil {
		log.Println("Database already initialized")
		return nil
	}

	// Log if database file will be created
	if _, err := os.Stat(config.DBLocation); errors.Is(err, os.ErrNotExist) {
		log.Println("Database file does not exist, will be created by SQLite on connect...")
	}

	// Open database connection
	db, err := sql.Open("sqlite3", connectionString(config))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection is working
	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("database not responding: %w", err)
	}

	// Run migrations
	if err = runMigrations(db, config.MigrationFiles); err != nil {
		db.Close()
		return fmt.Errorf("migration failed: %w", err)
	}

	// Store connection
	dbConnection = db
	config.Connection = db

	log.Println("Database initialized successfully")
	return nil
}

// GetDB returns the database connection.
// Returns nil if database is not initialized.
func GetDB() *sql.DB {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return dbConnection
}

// CloseDB safely closes the database connection.
// Safe to call multiple times; only closes once.
func CloseDB() {
	closeOnce.Do(func() {
		dbMutex.Lock()
		defer dbMutex.Unlock()

		if dbConnection == nil {
			log.Println("Database connection already closed or never opened")
			return
		}

		if err := dbConnection.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}

		dbConnection = nil
	})
}

// IsConnected checks if the database connection is active.
func IsConnected() bool {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	if dbConnection == nil {
		return false
	}

	return dbConnection.Ping() == nil
}

// runMigrations applies database schema migrations.
func runMigrations(db *sql.DB, migrationFiles string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationFiles),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to run")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// connectionString constructs a SQLite connection string.
func connectionString(config *models.DBConfig) string {
	return fmt.Sprintf(
		"%s?_auto_vacuum=%s&_foreign_keys=%t",
		config.DBLocation,
		config.Vacuum,
		config.ForeignKeys,
	)
}
