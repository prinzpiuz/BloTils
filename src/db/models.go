// Package db provides functionality for interacting with the application's database.
package db

import "time"

// DomainSettings represents the settings for a domain, including whether likes and comments are enabled,
// and the timestamp of the last update.
type DomainSettings struct {
	id        int
	likes     bool
	comments  bool
	timestamp time.Time
}

// Domain represents a domain and its associated settings.
// The id field is the unique identifier for the domain.
// The settings field contains the domain-specific settings.
// The domain field contains the domain name.
// The timestamp field contains the timestamp for when the domain was created or updated.
type Domain struct {
	id        int
	settings  DomainSettings
	domain    string
	timestamp time.Time
}

// IsEmpty returns true if the Domain is the zero value.
func (domain Domain) IsEmpty() bool {
	return domain == Domain{}
}

// LikesEnabled returns whether likes are enabled for the given Domain.
func (domain Domain) LikesEnabled() bool {
	return domain.settings.likes
}

// CommentsEnabled returns whether comments are enabled for the given Domain.
func (domain Domain) CommentsEnabled() bool {
	return domain.settings.comments
}

// Like represents a like for a domain and URI.
// The id field is the unique identifier for the like.
// The domain field is the domain the like is for.
// The uri field is the URI the like is for.
// The count field is the number of likes for the domain and URI.
// The timestamp field is the timestamp of when the like was created.
type Likes struct {
	id        int
	domain    Domain
	uri       string
	count     int
	timestamp time.Time
}

func (likes Likes) IsEmpty() bool {
	return likes == Likes{}
}

// LikedIPs represents a record of an IP address that has liked something, along with the count of likes and the timestamp of the last like.
type LikedIPs struct {
	id        int
	ip        string
	count     int
	timestamp time.Time
}
