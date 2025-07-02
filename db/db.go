package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// declared to access it outside this file
func Connect() error { //connection function returns errror
	var err error
	connect := "host=localhost port=5432 user=postgres dbname=tasks sslmode=disable"
	DB, err = sql.Open("postgres", connect)
	if err != nil {
		return fmt.Errorf("error opening database %s", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("unnable to connect to database %s", err)
	}
	//no error
	return nil
}
