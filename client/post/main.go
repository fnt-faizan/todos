package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	Id      string `json:"id"` //changed to string for uuid support
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

func main() {

	var todo Todo = Todo{Id: "ast", Title: "Hello World 2", Status: false, Deleted: false}
	dat, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error Marshalling")
	} else {
		res, err := http.Post("http://localhost:8080/todos/", "application/json; charset=utf-8", bytes.NewBuffer(dat))
		if err != nil {
			fmt.Println("Error making request", err)
		} else {
			defer res.Body.Close()
			var t Todo
			fmt.Println(res.Status)
			json.NewDecoder(res.Body).Decode(&t)
			fmt.Println(t)
		}
	}
}
