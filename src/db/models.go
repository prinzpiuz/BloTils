package db

import "time"

type Source struct {
	id        int
	source    string
	like      bool
	comment   bool
	timestamp time.Time
}

type LikedIP struct {
	id        int
	ip        string
	timestamp time.Time
}

type Like struct {
	id        int
	uri       string
	count     int
	timestamp time.Time
}
