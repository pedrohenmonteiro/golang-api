package database

import "database/sql"

type User struct {
}

func NewUser(db *sql.DB) *User {
	return &User{}
}
