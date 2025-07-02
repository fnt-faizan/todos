package main

import (
	"fmt"
	"net/http"
)

// testing delete
func main() {
	var id string
	fmt.Scan(&id)
	url := "http://localhost:8080/todos/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		return
	} else {
		cl := http.Client{}
		res, err := cl.Do(req)
		//fmt.Printf("Type of res is %T", res)
		if err != nil {
			return
		} else {
			fmt.Println((*res).Status)
		}
		defer res.Body.Close()
	}
}
