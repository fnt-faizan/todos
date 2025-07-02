package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	connect := "host=localhost port=5432 user=postgres password=Ttcmp2001## dbname=test sslmode=disable"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		fmt.Println("Error opening db")
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Println("Unable to connect to database", err)
	} else {
		fmt.Println("Connected to db")

		rows, err := db.Query("Select id, title from todo")
		if err != nil {
			fmt.Println("Error getting rows")
		} else {
			fmt.Printf("Type of rows is %T \n", rows)
			for rows.Next() {
				var id int
				var title string

				rows.Scan(&id, &title)
				fmt.Println(id, title)
			}
		}
	}
}
