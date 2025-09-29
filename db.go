package portfolio

import (
	"fmt"
	"log"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB(dbType string) *gorm.DB {
	var err error
	var db *gorm.DB

	switch dbType {
	case "postgres":
		db, err = connectPostgres()
	case "sqlite":
		db, err = connectSQLite()
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	
	err = db.AutoMigrate(
		&models.User{},
		&models.BlogPost{},
		&models.BlogComment{},
		&models.PostLike{},
		&models.CommentLike{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database models: %v", err)
	}

	log.Println("Connected to the database successfully")
	return db
}

func connectPostgres() (*gorm.DB, error) {
	host := internal.Env("DB_HOST")
    port := internal.Env("DB_PORT")
    user := internal.Env("DB_USER")
    password := internal.Env("DB_PASSWORD")
    dbname := internal.Env("DB_NAME")
    sslmode := internal.Env("DB_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("one or more required environment variables for Postgres are missing, required: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	
	return gorm.Open(postgres.Open(dns), &gorm.Config{})
}

func connectSQLite() (*gorm.DB, error) {
	dbPath := internal.Env("DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH environment variable is not set for SQLite")
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}