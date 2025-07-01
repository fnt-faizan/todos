package main

import (
	"fmt"
	"net/http"
)

// testing delete
func main() {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todos/1", nil)
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
