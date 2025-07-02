package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// This creates a server to explicitly handle todo requests
type Todo struct {
	Id      string `json:"id"` //changed to string for uuid support
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

// make a map of todo pointers to handle creates in memory
var todos = make(map[string]*Todo)

// declare a global variable for unique ids
//var id = 1 --> moved to UUIDs

// mutex for concurrency
var mutex = sync.RWMutex{}

// function to list all todos
func listalltodos(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()         //read lock
	defer mutex.RUnlock() //unlock
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
	mutex.Lock()
	defer mutex.Unlock()
	//get the data from request
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		//send a bad request status if data isn't valid
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	t.Id = uuid.New().String() //assigned a unique id different from the request id

	todos[t.Id] = &t //store to a map
	//send a created respose
	w.WriteHeader(http.StatusCreated)
	// send back the created data
	json.NewEncoder(w).Encode(t)
}

// this function return a specific todo by id
func listtodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	mutex.RLock()
	defer mutex.RUnlock()
	id := r.URL.Path[len("/todos/"):] //returns string after the path
	t, ok := todos[id]                //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
}
func deltodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	mutex.Lock()
	defer mutex.Unlock()
	id := r.URL.Path[len("/todos/"):] //returns string after the path
	t, ok := todos[id]                //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	t.Deleted = true
	w.WriteHeader(http.StatusNoContent) //deleted
}
func updatetodo(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	// get the id
	id := r.URL.Path[len("/todos/"):] //returns string after the path
	t, ok := todos[id]                //return a bool if val at id exists
	if !ok || t.Deleted {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}
	err := json.NewDecoder(r.Body).Decode(todos[id])
	todos[id].Id = id
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)

}
func main() {
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
