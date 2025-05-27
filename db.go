package portfolio

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string, migrationsDir string, seedsDir string) {
	db.DataBase = connectDB(dbPath)

	runMigrations(migrationsDir, db.DataBase)
	seedUsersIfEmpty(seedsDir, db.DataBase)
}

func connectDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func runMigrations(migrationsDir string, DB *gorm.DB) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		panic("failed to read migrations directory: " + err.Error())
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		sqlPath := filepath.Join(migrationsDir, file.Name())
		content, err := ioutil.ReadFile(sqlPath)
		if err != nil {
			panic("failed to read migration file: " + err.Error())
		}
		if err := DB.Exec(string(content)).Error; err != nil {
			panic("failed to execute migration: " + err.Error())
		}
	}
}

func seedUsersIfEmpty(seedsDir string, DB *gorm.DB) {
    var count int64
    if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
        panic("failed to count users: " + err.Error())
    }
    if count > 0 {
        return // Table is not empty, skip seeding
    }

    files, err := ioutil.ReadDir(seedsDir)
    if err != nil {
        panic("failed to read seeds directory: " + err.Error())
    }

    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
            continue
        }
        jsonPath := filepath.Join(seedsDir, file.Name())
        content, err := ioutil.ReadFile(jsonPath)
        if err != nil {
            panic("failed to read seed file: " + err.Error())
        }
        var users []models.User
        if err := json.Unmarshal(content, &users); err != nil {
            panic("failed to unmarshal seed file: " + err.Error())
        }
		// hash the passwords before inserting
		for i := range users {
			if users[i].Password != "" {	
				hashedPassword := internal.Encrypt(users[i].Password)
				users[i].Password = hashedPassword
			}
		}
        if len(users) > 0 {
            if err := DB.Create(&users).Error; err != nil {
                panic("failed to insert seed data: " + err.Error())
            }
        }
    }
}