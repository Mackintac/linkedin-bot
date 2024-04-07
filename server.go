package main

import (
	"fmt"
	"log"
	"net/http"
)

func Server() error {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	fmt.Println("Server going live")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server dead")
		log.Fatal("Error with ListenAndServe", err)
	}

	return nil
}
