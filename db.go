package portfolio

import (
	"fmt"
	"log"
	"net/url"

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
		&models.Portfolio{},
		&models.Skill{},
		&models.SoftSkill{},
		&models.Language{},
		&models.Education{},
		&models.EducationTask{},
		&models.Hobby{},
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
		return nil, fmt.Errorf("One or more required PostgreSQL config settings are not set")
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   dbname,
	}
	query := dsn.Query()
	query.Set("sslmode", sslmode)
	dsn.RawQuery = query.Encode()

	return gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
}

func connectSQLite() (*gorm.DB, error) {
	dbPath := config.AppConfig.Database.Path
	if dbPath == "" {
		return nil, fmt.Errorf("Database path is not set for SQLite")
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
