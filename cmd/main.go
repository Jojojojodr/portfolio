package main

import (
	"log"
	"os"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/seed"
)

func main() {
	log.Println("Starting Application")

	// Set up database connection and migrations
	db.DataBase = portfolio.ConnectDB()

	// Seed the database with initial data
	if internal.IsDatabaseEmpty(db.DataBase) {
		log.Println("Database is empty, seeding with initial data")
		if err := seed.SeedDatabase(db.DataBase); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	} else {
		log.Println("Database already contains data, skipping seeding")
	}

	// Run the server
	port := internal.Env("PORT")
	if port == "" {
		log.Println("PORT not set")
		os.Exit(1)
	}

	portfolio.RunServer(port)
}
