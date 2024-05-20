package db

import "time"

type DomainSettings struct {
	id        int
	likes     bool
	comments  bool
	timestamp time.Time
}

type Domain struct {
	id        int
	settings  DomainSettings
	domain    string
	timestamp time.Time
}

type Like struct {
	id        int
	domain    Domain
	uri       string
	count     int
	timestamp time.Time
}

type LikedIPs struct {
	id        int
	ip        string
	count     int
	timestamp time.Time
}
