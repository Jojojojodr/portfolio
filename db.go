package portfolio

import (
	"fmt"
	"log"

	"github.com/Jojojojodr/portfolio/config"
	"github.com/Jojojojodr/portfolio/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Data *Database

type Database struct {
	db *gorm.DB
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func (d *Database) SetDB(db *gorm.DB) {
	d.db = db
	Data = d
}

func (d *Database) Connect(dbType string) error {
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
		return fmt.Errorf("Failed to connect to the database: %v", err)
	}

	d.SetDB(db)

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
	return nil
}

func NewDatabase() *Database {
	return &Database{db: nil}
}

func connectPostgres() (*gorm.DB, error) {
	host := config.AppConfig.Database.Host
	port := config.AppConfig.Database.Port
	user := config.AppConfig.Database.Username
	password := config.AppConfig.Database.Password
	dbname := config.AppConfig.Database.Name
	sslmode := config.AppConfig.Database.SSLMode

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
	dbPath := config.AppConfig.Database.Path
	if dbPath == "" {
		return nil, fmt.Errorf("DB_PATH environment variable is not set for SQLite")
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
