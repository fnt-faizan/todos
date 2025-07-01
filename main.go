package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// declare a global variable for unique ids
var id = 1

// function to list all todos
func listalltodos(w http.ResponseWriter, r *http.Request) {
	var td []*Todo
	for _, val := range todos {
		if !val.Deleted {
			td = append(td, val)
		}
	}
	if len(td) == 0 {
		http.Error(w, "No todos to show", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(td)
}

// function to create a todo
func createtodo(w http.ResponseWriter, r *http.Request) {
	//get the data from request
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		//send a bad request status if data isn't valid
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	t.Id = id //assigned a unique id different from the request id
	id++
	todos[t.Id] = &t //store to a map
	//send a created respose
	w.WriteHeader(http.StatusCreated)
	// send back the created data
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
	t, ok := todos[id] //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}
func deltodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	strpath := r.URL.Path[len("/todos/"):] //returns string after the path
	id, err := strconv.Atoi(strpath)
	if err != nil {
		http.Error(w, "Inavalid todo ID", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Inavalid todo ID", http.StatusNotFound)
		return
	}
	t, ok := todos[id] //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	t.Deleted = true
	w.WriteHeader(http.StatusNoContent) //deleted
}
func updatetodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	strpath := r.URL.Path[len("/todos/"):] //returns string after the path
	id, err := strconv.Atoi(strpath)
	if err != nil {
		http.Error(w, "Inavalid todo ID", http.StatusNotFound)
		return
	}
	t, ok := todos[id] //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	err = json.NewDecoder(r.Body).Decode(todos[id])
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)

}
func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			listalltodos(w, r)
		} else if r.Method == "POST" {
			createtodo(w, r)
		} else {
			http.Error(w, "Method is not allwed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		//id , _ := strconv.Atoi(URL.Path[len("/todos/"):1])
		if r.Method == "GET" {
			listtodo(w, r)
		} else if r.Method == "POST" {
			createtodo(w, r)
		} else if r.Method == "PUT" {
			updatetodo(w, r)
		} else if r.Method == "DELETE" {
			deltodo(w, r)
		} else {
			http.Error(w, "Method is not allwed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server started on port 8080")
	http.ListenAndServe("localhost:8080", nil)
}
