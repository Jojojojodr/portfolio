package models

import (
	"log"

	"github.com/Jojojojodr/portfolio/internal/db"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
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

func GetUsers() []User {
	var users []User
	result := db.DataBase.Find(&users)
	if result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
	}
	return users
}

func GetUserByID(id uint) (*User, error) {
    var user User
    result := db.DataBase.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func GetUserByName(name string) (*User, error) {
	var user User
	result := db.DataBase.Where("name = ?", name).First(&user)
	if result.Error != nil {
		log.Printf("Error fetching user by name %s: %v", name, result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
    var user User
    result := db.DataBase.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func UpdateUser(user *User) error {
    result := db.DataBase.Save(user)
    return result.Error
}

func CountAdminUsers() (int64, error) {
    var count int64
    result := db.DataBase.Model(&User{}).Where("is_admin = ?", true).Count(&count)
    return count, result.Error
}