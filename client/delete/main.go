package main

import (
	"fmt"
	"net/http"
)

// testing delete
func main() {

	hosturl := "http://localhost:8080/todos/" + "a34ff89b-95f8-4f7d-82b6-470f892d1a21"
	req, err := http.NewRequest(http.MethodDelete, hosturl, nil)
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
