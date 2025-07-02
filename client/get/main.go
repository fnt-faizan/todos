package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	Id      string `json:"id"` //changed to string for uuid support
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

// This program fetches a TODO item from a public API and prints it to the console.
func main() {
	var id string
	fmt.Scan(&id)
	url := "http://localhost:8080/todos/" + id
	res, err := http.Get(url)
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
