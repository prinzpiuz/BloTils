package models

import "database/sql"

type DBConfig struct {
	DBLocation      string
	Vacuum          string
	ForeignKeys     bool
	Connection      *sql.DB
	DBBaseDirectory string
}

func (dbConfig DBConfig) IsEmpty() bool {
	return dbConfig == DBConfig{}
}
