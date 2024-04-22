package main

import (
	"dev/linkedIn/components/handlers"

	projectUtil "dev/linkedIn/util"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

var projectConfig projectUtil.TProjectConfig = projectUtil.InitProjectConfig()

func main() {

	handlers.InitHandlers()
	if err := Server(); err != nil {
		log.Fatal("Error starting server")
		return
	}
}
