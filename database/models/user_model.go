package models

import (
	"time"
	"log"

	"gorm.io/gorm"
)

type User struct {
	ID       	uint   `json:"id" gorm:"primaryKey"`
	Name     	string `json:"name"`
	Email    	string `json:"email"`
	Password 	string `json:"-"`
	IsAdmin  	bool   `json:"is_admin"`
	CreatedAt 	time.Time `json:"created_at"`
}

func (user *User) IsAuthenticated(tx *gorm.DB) bool {
	if user.ID == 0 {
		return false
	}

	result := tx.First(user, user.ID)
	if result.Error != nil {
		log.Printf("Error fetching user with ID %d: %v", user.ID, result.Error)
		return false
	}

	return user.ID > 0 && user.Name != "" && user.Email != ""
}