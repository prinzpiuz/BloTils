// Package db provides utility functions for interacting with the database.
// this page contains the utility functions for db package
package db

import (
	"database/sql"
	"log"
)

// GetDomain retrieves a Domain from the database by the given domain_name.
// If the domain is not found, it returns an empty Domain.
func GetDomain(db *sql.DB, domain_name string) Domain {
	var domain Domain
	err := db.QueryRow(getDomainQuery, domain_name).Scan(&domain.id,
		&domain.settings.id,
		&domain.domain,
		&domain.timestamp,
		&domain.settings.id,
		&domain.settings.likes,
		&domain.settings.comments,
		&domain.settings.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Domain %s, Not Found", domain_name)
			return Domain{}
		}
		log.Printf("Error Getting Domain %s: %v", domain_name, err)
		return Domain{}
	}
	return domain
}

func GetLikes(db *sql.DB, domain_name string, page string) Likes {
	var likes Likes
	err := db.QueryRow(getLikesQuery, domain_name, page).Scan(&likes.id,
		&likes.uri,
		&likes.count,
		&likes.domain.id,
		&likes.timestamp,
		&likes.domain.id,
		&likes.domain.settings.id,
		&likes.domain.domain,
		&likes.domain.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: Likes On %s, Not Found In %s", page, domain_name)
			return Likes{}
		}
		log.Printf("Error Getting Likes On %s For %s: %v", page, domain_name, err)
		return Likes{}
	}
	return likes
}
