package main

import "fmt"

func main(){
	srv, err := NewServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	srv.Listen()
}
