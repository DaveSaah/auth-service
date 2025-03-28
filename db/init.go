// package db defines functions related to the database
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// Init initialises a new database
func Init() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return nil, err
	}

	return conn, err
}
