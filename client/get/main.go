package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// This program fetches a TODO item from a public API and prints it to the console.
func main() {
	res, err := http.Get("http://localhost:8080/todos")
	if err != nil {
		fmt.Println("Error fetching data:", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.Status)
	} else if res.StatusCode == http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			var todo Todo
			json.Unmarshal(data, &todo)
			fmt.Println(todo)
		}
	}
	// var todo Todo
	// err = json.NewDecoder(res.Body).Decode(&todo)
	// if err != nil {
	// 	fmt.Println("Error decoding JSON:", err)
	// } else {
	// 	fmt.Println(todo)
	// 	fmt.Printf("User ID: %d, ID: %d, Title: %s, Completed: %t\n", todo.UserID, todo.ID, todo.Title, todo.Completed)
	// }
}
