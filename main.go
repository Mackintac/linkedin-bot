package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading information from .env")
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	}
	http.HandleFunc("/", handler)
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	fmt.Println("Server live")
	log.Fatal(server.ListenAndServe())
}
