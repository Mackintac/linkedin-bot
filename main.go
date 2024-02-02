package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	if err := envInit(); err != nil {
		log.Fatal("Error loading information from .env")
		return
	}

	if err := helloHandler(); err != nil {
		log.Fatal("Error loading Hello handler")
		return
	}

	fmt.Println("Server going live")
	if err := startServer(); err != nil {
		log.Fatal("Error starting server")
		return
	}
}

func envInit() error {
	// .env initialization
	godotenv.Load(".env")
	return nil
}

func helloHandler() error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
		fmt.Println("Server Pinged!")
	}
	http.HandleFunc("/", handler)
	return nil
}

func startServer() error {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server dead")
		log.Fatal("Error with ListenAndServe")
	}
	return nil
}
