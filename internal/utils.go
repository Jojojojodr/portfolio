package internal

import (
	"fmt"
	"os"

	"github.com/Jojojojodr/portfolio/internal/db/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Env(key string) string {
	return os.Getenv(key)
}

func Encrypt(key string) string {
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error encrypting key: %v\n", err)
		return ""
	}
	return string(hashedKey)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsDatabaseEmpty(database *gorm.DB) bool {
    var userCount int64
    
    database.Model(&models.User{}).Count(&userCount)
    
    return userCount == 0
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
}