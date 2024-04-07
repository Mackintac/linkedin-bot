package main

import (
	"dev/linkedIn/components/handlers"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	handlers.InitHandlers()

	if err := Server(); err != nil {
		log.Fatal("Error starting server")
		return
	}

}
