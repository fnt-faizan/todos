package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todos/db"
	"todos/models"

	"github.com/gorilla/mux"
)

const (
	cacheListKey = "todos:all"
	itemKeyFmt   = "todo:%d"
)

// helper functions for JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, msg string, status int) {
	jsonResponse(w, map[string]string{"error": msg}, status)
}

//helper end

// Handlers for the Todo API
func listAllTodos(w http.ResponseWriter, r *http.Request) {
	if data, err := db.RDB.Get(db.RCtx, cacheListKey).Bytes(); err == nil {
		var todos []*models.Todo
		json.Unmarshal(data, &todos)
		jsonResponse(w, todos, http.StatusOK)
		return
	}

	todos, err := db.GetAllTodos(r.Context(), db.DB)
	if err != nil {
		errorResponse(w, "database error", http.StatusInternalServerError)
		return
	}
	if len(todos) == 0 {
		errorResponse(w, "no todos to be shown", http.StatusNotFound)
		return
	}

	buf, _ := json.Marshal(todos)
	db.RDB.Set(db.RCtx, cacheListKey, buf, 5*time.Minute)
	jsonResponse(w, todos, http.StatusOK)
}

func listTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	itemKey := fmt.Sprintf(itemKeyFmt, id)
	if data, err := db.RDB.Get(db.RCtx, itemKey).Bytes(); err == nil {
		var t models.Todo
		json.Unmarshal(data, &t)
		jsonResponse(w, t, http.StatusOK)
		return
	}

	todo, err := db.GetTodoByID(r.Context(), db.DB, id)
	if err != nil {
		errorResponse(w, "todo not found or deleted", http.StatusNotFound)
		return
	}

	buf, _ := json.Marshal(todo)
	db.RDB.Set(db.RCtx, itemKey, buf, 5*time.Minute)
	jsonResponse(w, todo, http.StatusOK)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	success, err := db.DeleteTodo(r.Context(), db.DB, id)
	if err != nil {
		errorResponse(w, "unable to delete todo", http.StatusInternalServerError)
		return
	}
	if !success {
		errorResponse(w, "todo not found", http.StatusNotFound)
		return
	}
	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, id))
	jsonResponse(w, map[string]string{"message": "todo deleted"}, http.StatusOK)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, "invalid todo ID", http.StatusBadRequest)
		return
	}

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		errorResponse(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	t.Id = id

	success, err := db.UpdateTodo(r.Context(), db.DB, &t)
	if err != nil {
		errorResponse(w, "update failed due to database error", http.StatusInternalServerError)
		return
	}
	if !success {
		errorResponse(w, "todo not found in records", http.StatusNotFound)
		return
	}
	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, id))
	jsonResponse(w, t, http.StatusOK)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		errorResponse(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := db.InsertTodo(r.Context(), db.DB, &t); err != nil {
		errorResponse(w, "insert into database failed", http.StatusInternalServerError)
		return
	}
	db.RDB.Del(db.RCtx, cacheListKey)
	db.RDB.Del(db.RCtx, fmt.Sprintf(itemKeyFmt, t.Id))
	jsonResponse(w, t, http.StatusCreated)
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("postgres connection failed:", err)
		return
	}
	defer db.DB.Close()

	if !db.ConnectRedis() {
		log.Fatal("redis connection failed")
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
