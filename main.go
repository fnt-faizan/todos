package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todos/db"

	"github.com/redis/go-redis/v9"
)

// Todo struct for JSON mapping
type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

var (
	DB   *sql.DB       = db.DB
	RDB  *redis.Client = db.RDB
	RCtx               = db.RCtx
)

const (
	cacheListKey = "todos:all" // Redis key for full todos list
	itemKeyFmt   = "todo:%d"   // Template for per-item cache key //can use sprintf to format
)

// listalltodos handles GET /todos; uses cache-aside strategy
func listalltodos(w http.ResponseWriter, r *http.Request) {
	// Try reading from cache
	if data, err := RDB.Get(RCtx, cacheListKey).Bytes(); err == nil {
		var td []*Todo
		json.Unmarshal(data, &td)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(td)
		return
	}

	// On cache miss, query DB
	var td []*Todo
	rows, err := DB.Query("SELECT id, title, status, deleted FROM todos WHERE deleted = false ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		t := new(Todo)
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

	// Cache the result
	buf, _ := json.Marshal(td)
	RDB.Set(RCtx, cacheListKey, buf, 5*time.Minute)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(td)
}

// listtodo handles GET /todos/{id}; caches individual items
func listtodo(w http.ResponseWriter, r *http.Request) {
	// get the id
	strpath := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(strpath)
	if err != nil {
		http.Error(w, "Inavalid todo ID", http.StatusNotFound)
		return
	}

	itemKey := fmt.Sprintf(itemKeyFmt, id)
	if data, err := RDB.Get(RCtx, itemKey).Bytes(); err == nil {
		var t Todo
		json.Unmarshal(data, &t)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)
		return
	}

	// cache miss -> load from DB
	var t Todo
	t.Id = id
	err = DB.QueryRow("SELECT title, status, deleted FROM todos WHERE id = $1", id).
		Scan(&t.Title, &t.Status, &t.Deleted)
	if err != nil {
		http.Error(w, "Item doesn't exist or may have been deleted", http.StatusNotFound)
		return
	}

	// cache the individual item
	buf, _ := json.Marshal(t)
	RDB.Set(RCtx, itemKey, buf, 5*time.Minute)

	json.NewEncoder(w).Encode(t)
}

// deltodo handles DELETE and invalidates cache
func deltodo(w http.ResponseWriter, r *http.Request) {
	strID := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// check if the id exists & delete
	res, err := DB.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Unable to delete todo", http.StatusInternalServerError)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// invalidate full list and individual item cache
	RDB.Del(RCtx, cacheListKey)
	RDB.Del(RCtx, fmt.Sprintf(itemKeyFmt, id))

	w.WriteHeader(http.StatusNoContent)
}

func updatetodo(w http.ResponseWriter, r *http.Request) {
	strID := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	res, err := DB.Exec(
		"UPDATE todos SET title = $1, status = $2, deleted = $3 WHERE id = $4",
		t.Title, t.Status, t.Deleted, id,
	)
	if err != nil {
		http.Error(w, "Unable to update todo", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	// invalidate caches
	RDB.Del(RCtx, cacheListKey)
	RDB.Del(RCtx, fmt.Sprintf(itemKeyFmt, id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func createtodo(w http.ResponseWriter, r *http.Request) {
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err := DB.QueryRow(`
        INSERT INTO todos (title, status, deleted)
        VALUES ($1, $2, $3)
        RETURNING id
    `, t.Title, t.Status, t.Deleted).Scan(&t.Id)
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		return
	}

	// invalidate caches
	RDB.Del(RCtx, cacheListKey)
	RDB.Del(RCtx, fmt.Sprintf(itemKeyFmt, t.Id))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func main() {
	// Create a connection to postgres
	if db.Connect() != nil {
		fmt.Println("Connection to database failed")
		return
	}
	DB = db.DB

	// Connect to redis
	if !db.ConnectRedis() {
		fmt.Println("Connection to redis failed")
		return
	}
	RDB = db.RDB
	RCtx = db.RCtx

	defer DB.Close()

	// Register endpoints
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listalltodos(w, r)
		case "POST":
			createtodo(w, r)
		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listtodo(w, r)
		case "PUT":
			updatetodo(w, r)
		case "DELETE":
			deltodo(w, r)
		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port 8080")
	http.ListenAndServe("localhost:8080", nil)
}
