package main

import (
	"log"
	"metrika/cmd/api"
	"metrika/internal/database"
)

func main() {
	db, err := database.PostgresConnect()
	if err != nil {
		log.Fatalln(err)
	}

	server := api.NewApiServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}

}
