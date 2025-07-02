package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

func main() {

	var todo Todo = Todo{Id: 1, Title: "Pushed via a script into DB haha", Status: false, Deleted: false}
	dat, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error Marshalling")
	} else {
		res, err := http.Post("http://localhost:8080/todos", "application/json; charset=utf-8", bytes.NewBuffer(dat))
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
