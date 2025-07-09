package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todos/db"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

const (
	cacheListKey = "todos:all"
	itemKeyFmt   = "todo:%d"
)

func listAllTodos(w http.ResponseWriter, r *http.Request) {
	if data, err := db.RDB.Get(db.RCtx, cacheListKey).Bytes(); err == nil {
		var todos []*Todo
		json.Unmarshal(data, &todos)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todos)
		return
	}

	rows, err := db.DB.Query("SELECT id, title, status, deleted FROM todos WHERE deleted = false ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		t := new(Todo)
		if err := rows.Scan(&t.Id, &t.Title, &t.Status, &t.Deleted); err != nil {
			http.Error(w, "Error scanning record", http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}

	if len(todos) == 0 {
		http.Error(w, "No todos to show", http.StatusNotFound)
		return
	}

	buf, _ := json.Marshal(todos)
	db.RDB.Set(db.RCtx, cacheListKey, buf, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func listTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	itemKey := fmt.Sprintf(itemKeyFmt, id)
	if data, err := db.RDB.Get(db.RCtx, itemKey).Bytes(); err == nil {
		var t Todo
		json.Unmarshal(data, &t)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)
		return
	}

	var t Todo
	t.Id = id
	err = db.DB.QueryRow("SELECT title, status, deleted FROM todos WHERE id = $1 AND deleted = false", id).
		Scan(&t.Title, &t.Status, &t.Deleted)
	if err != nil {
		http.Error(w, "Item not found or deleted", http.StatusNotFound)
		return
	}

	buf, _ := json.Marshal(t)
	db.RDB.Set(db.RCtx, itemKey, buf, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Soft delete by setting 'deleted' = true
	res, err := db.DB.Exec("UPDATE todos SET deleted = true WHERE id = $1 where deleted = false", id)
	if err != nil {
		http.Error(w, "Unable to delete todo", http.StatusInternalServerError)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// Invalidate cache
	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, id))

	w.WriteHeader(http.StatusNoContent)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec(
		"UPDATE todos SET title = $1, status = $2, deleted = $3 WHERE id = $4 and deleted = false",
		t.Title, t.Status, t.Deleted, id,
	)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow(`
        INSERT INTO todos (title, status, deleted)
        VALUES ($1, $2, $3)
        RETURNING id
    `, t.Title, t.Status, t.Deleted).Scan(&t.Id)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, t.Id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Postgres connection failed:", err)
		return
	}
	defer db.DB.Close()

	if !db.ConnectRedis() {
		log.Fatal("Redis connection failed")
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/todos", listAllTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", listTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}
