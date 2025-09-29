package main

import (
	"flag"
	"log"
	"os"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/seed"
)

func main() {
	var port = flag.String("p", "", "Port to run the server on (e.g., 8080)")
	var dbType = flag.String("d", "", "Database type: sqlite or postgres")
	var token = flag.String("t", "", "Secret token for JWT authentication")
	flag.Parse()

	log.Println("Starting Application")

	secretToken := *token
	if secretToken == "" {
		secretToken = internal.Env("SECRET_TOKEN")
		if secretToken == "" {
			log.Println("Secret token not specified via -t flag or SECRET_TOKEN environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite -t your_secret_token")
			os.Exit(1)
		}
	}
	internal.SetSecretToken(secretToken)

	selectedDBType := *dbType
	if selectedDBType == "" {
		selectedDBType = internal.Env("DB_TYPE")
		if selectedDBType == "" {
			log.Println("Database type not specified via -d flag or DB_TYPE environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite")
			os.Exit(1)
		}
	}

	db.DataBase = portfolio.ConnectDB(selectedDBType)

	if internal.IsDatabaseEmpty(db.DataBase) {
		log.Println("Database is empty, seeding with initial data")
		if err := seed.SeedDatabase(db.DataBase); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	} else {
		log.Println("Database already contains data, skipping seeding")
	}

	selectedPort := *port
	if selectedPort == "" {
		selectedPort = internal.Env("PORT")
		if selectedPort == "" {
			log.Println("Port not specified via -p flag or PORT environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite")
			os.Exit(1)
		}
	}

	log.Printf("Starting server on port %s with database type %s", selectedPort, selectedDBType)
	portfolio.RunServer(selectedPort)
}
