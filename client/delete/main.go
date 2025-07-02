package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// testing delete
func main() {
	for i := 2; i < 20; i++ {
		hosturl := "http://localhost:8080/todos/" + strconv.Itoa(i)
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
}
