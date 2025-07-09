package main

import (
	"fmt"
	"net/http"
)

// testing delete
func main() {

	hosturl := "http://localhost:8080/todos/" + "1" // replace 1 with the ID of the todo you want to delete
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
