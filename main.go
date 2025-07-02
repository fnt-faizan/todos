package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todos/db"
)

// This creates a server to explicitly handle todo requests
type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

// make a map of todo pointers to handle creates in memory
var todos = make(map[int]*Todo)

// -> not needed now stored in db

var DB *sql.DB

// declare a global variable for unique ids
//var id = 2 -> not needed postgre auto increments

// function to list all todos
func listalltodos(w http.ResponseWriter, r *http.Request) {
	var td []*Todo
	rows, err := DB.Query("SELECT id, title,status, deleted from todos where deleted = false order by id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := new(Todo) //returns a pointer of type todo
		if err := rows.Scan(&t.Id, &t.Title, &t.Status, &t.Deleted); err != nil {
			http.Error(w, "Error fetching record", http.StatusInternalServerError)
			return
		}
		td = append(td, t)
	}
	if len(td) == 0 {
		http.Error(w, "No todos to show", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(td)
}

func createtodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Use QueryRow with RETURNING to get the auto-generated ID
	query := `
      INSERT INTO todos (title, status, deleted)
      VALUES ($1, $2, $3)
      RETURNING id
    `
	err := DB.QueryRow(query, t.Title, t.Status, t.Deleted).Scan(&t.Id)
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

// this function return a specific todo by id
func listtodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	strpath := r.URL.Path[len("/todos/"):] //returns string after the path
	id, err := strconv.Atoi(strpath)
	if err != nil {
		http.Error(w, "Inavalid todo ID", http.StatusNotFound)
		return
	}
	row := DB.QueryRow("SELECT title, status, deleted from todos where id = $1", id)
	var t Todo
	t.Id = id
	err = row.Scan(&t.Title, &t.Status, &t.Deleted)
	if err != nil {
		//zero rows fetched
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}
func deltodo(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID
	strID := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	//chheck if the id exists
	res, err := DB.Exec("SELECT id from todos where id = $1", id)
	if err != nil {
		http.Error(w, "Falied to locate todo", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	// Run a delete query
	res, err = DB.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Unable to delete todo", http.StatusInternalServerError)
		return
	}

	// Check if a row was actually deleted
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func updatetodo(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID
	strID := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Parse JSON body into a struct
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Run the UPDATE query
	res, err := DB.Exec(
		"UPDATE todos SET title = $1, status = $2, deleted = $3 WHERE id = $4",
		t.Title, t.Status, t.Deleted, id,
	)
	if err != nil {
		http.Error(w, "Unable to update todo", http.StatusInternalServerError)
		return
	}

	// Ensure a row was updated
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Return the updated resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func main() {
	//Create a connection to postgres
	if db.Connect() != nil {
		fmt.Println("Connection to database failed")
		return
	}
	DB = db.DB
	defer DB.Close() //close connection
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listalltodos(w, r)
		case "POST":
			createtodo(w, r)
		default:
			http.Error(w, "Method is not allwed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		//id , _ := strconv.Atoi(URL.Path[len("/todos/"):1])
		switch r.Method {
		case "GET":
			listtodo(w, r)
		case "POST":
			createtodo(w, r)
		case "PUT":
			updatetodo(w, r)
		case "DELETE":
			deltodo(w, r)
		default:
			http.Error(w, "Method is not allwed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server started on port 8080")
	http.ListenAndServe("localhost:8080", nil)
}
