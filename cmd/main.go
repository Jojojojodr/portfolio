package main

import (
	"log"
	"os"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
)

func main() {
	log.Println("Starting Application")

	// Set up database connection and migrations
	db.DataBase = portfolio.ConnectDB("./database/sqlite.db")

	// Run the server
	port := internal.Env("PORT")
	if port == "" {
		log.Println("PORT not set")
		os.Exit(1)
	}

	portfolio.RunServer(port)
}
