package main

import (
	"log"
	"os"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/internal"
)

func main() {
	log.Println("Starting Application")

	// Set up database connection and migrations
	portfolio.InitDB("./database/sqlite.db", "./database/migrations", "./database/seeds")

	// Run the server
	port := internal.Env("PORT")
	if port == "" {
		log.Println("PORT not set")
		os.Exit(1)
	}

	portfolio.RunServer(port)
}
