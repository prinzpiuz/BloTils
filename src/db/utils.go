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
	err := db.QueryRow(getDomain, domain_name).Scan(&domain.ID,
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
	err := db.QueryRow(getLikes, domain_name, page).Scan(
		&likes.id,
		&likes.URI,
		&likes.domain_id,
		&likes.Count,
		&likes.Domain.ID,
		&likes.Domain.settings.id,
		&likes.Domain.domain,
		&likes.Domain.timestamp)
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

func GetLikedIP(db *sql.DB, domain string, page string, ip string) LikedIPs {
	var likedIP LikedIPs
	err := db.QueryRow(getIPlikedOrNot, ip, domain, page).Scan(&likedIP.id,
		&likedIP.IP,
		&likedIP.Count,
		&likedIP.Domain,
		&likedIP.Path,
		&likedIP.timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("DB: IP %s, Not Found Liked in %s for %s", ip, domain, page)
			return LikedIPs{}
		}
		log.Printf("Error Getting Likes On %s For %s: %v", ip, page, err)
		return LikedIPs{}
	}
	return likedIP
}

func UpdateIPLikeCount(db *sql.DB, domain string, path string, ip string) {
	_, err := db.Exec(updateOrInsertLikedIP, domain, path, ip)
	if err != nil {
		log.Printf("Error Updating IP Like Count: %v", err)
	}

}

func UpdateLikeCount(db *sql.DB, page string, doamin_id int) error {
	_, err := db.Exec(updateLike, page, doamin_id)
	return err
}
