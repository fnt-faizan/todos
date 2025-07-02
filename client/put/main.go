package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Todo struct {
	Id      string `json:"id"` //changed to string for uuid support
	Title   string `json:"title"`
	Status  bool   `json:"status"`
	Deleted bool   `json:"deleted"`
}

func main() {
	var id string
	fmt.Scan(&id)
	url := "http://localhost:8080/todos/" + id
	var todo Todo = Todo{Id: id, Title: "Modified", Status: false, Deleted: false}
	dat, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error Marshalling")
		return
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(dat))
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
