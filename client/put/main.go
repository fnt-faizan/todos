package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Todo struct {
	Id      int    `json:"userId"`
	Status  bool   `json:"status"`
	Title   string `json:"title"`
	Deleted bool   `json:"deleted"`
}

func main() {
	var todo Todo = Todo{Id: 1, Title: "Chnaged", Status: false, Deleted: false}
	dat, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error Marshalling")
		return
	}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/todos/a34ff89b-95f8-4f7d-82b6-470f892d1a21", bytes.NewBuffer(dat))
	if err != nil {
		fmt.Println("Request creation failed")
		return
	} else {
		//created request
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Request failed")
			return
		}
		defer res.Body.Close()
		fmt.Println(res.Status)
		dat, _ := io.ReadAll(res.Body)
		fmt.Println(string(dat))

	}
}
