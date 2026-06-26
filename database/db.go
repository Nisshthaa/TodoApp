package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Todo *sqlx.DB

// Create connection
func OpenConnection(host, port, databaseName, user, password string) error {

	connStr := fmt.Sprintf(

		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		databaseName,
	)

	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		return err
	}

	Todo = db
	return nil
}

// db close
