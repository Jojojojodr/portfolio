package internal

import (
	"fmt"

	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SecretToken string

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

func SetSecretToken(token string) {
	SecretToken = token
}