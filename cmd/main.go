package main

import (
	"flag"
	"log"
	"os"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/config"
	"github.com/Jojojojodr/portfolio/database/seed"
	"github.com/Jojojojodr/portfolio/routers"
)

func main() {
	config.LoadConfig()

	var port = flag.String("p", "", "Port to run the server on (e.g., 8080)")
	var dbType = flag.String("d", "", "Database type: sqlite or postgres")
	var token = flag.String("t", "", "Secret token for JWT authentication")
	flag.Parse()

	log.Println("Starting Application")

	secretToken := *token
	if secretToken == "" {
		secretToken = config.AppConfig.Server.JWTSecret
		if secretToken == "" {
			log.Println("Secret token not specified via -t flag or SECRET_TOKEN environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite -t your_secret_token")
			os.Exit(1)
		}
	}
	portfolio.SetSecretToken(secretToken)

	selectedDBType := *dbType
	if selectedDBType == "" {
		selectedDBType = config.AppConfig.Database.Type
		if selectedDBType == "" {
			log.Println("Database type not specified via -d flag or DB_TYPE environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite")
			os.Exit(1)
		}
	}

	database := portfolio.NewDatabase()
	if err := database.Connect(selectedDBType); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if portfolio.IsDatabaseEmpty(portfolio.Data.GetDB()) {
		log.Println("Database is empty, seeding with initial data")
		if err := seed.SeedDatabase(portfolio.Data.GetDB()); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	} else {
		log.Println("Database already contains data, skipping seeding")
	}

	selectedPort := *port
	if selectedPort == "" {
		selectedPort = config.AppConfig.Server.Port
		if selectedPort == "" {
			log.Println("Port not specified via -p flag or PORT environment variable")
			log.Println("Usage: ./app -p 8080 -d sqlite")
			os.Exit(1)
		}
	}

	log.Printf("Starting server on port %s with database type %s", selectedPort, selectedDBType)
	server := portfolio.NewServer(selectedPort)

	routers.FrontendRouter(server.Engine)
	routers.V1Router(server.Engine)
	routers.HandleRouter(server.Engine)

	server.Start()
}
