package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUsers() []models.User {
	var users []models.User
	result := portfolio.Data.GetDB().Find(&users)
	if result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
	}
	return users
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := portfolio.Data.GetDB().First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByName(name string) (*models.User, error) {
	var user models.User
	result := portfolio.Data.GetDB().Where("name = ?", name).First(&user)
	if result.Error != nil {
		log.Printf("Error fetching user by name %s: %v", name, result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := portfolio.Data.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	result := portfolio.Data.GetDB().Save(user)
	return result.Error
}

func CountAdminUsers() (int64, error) {
	var count int64
	result := portfolio.Data.GetDB().Model(&models.User{}).Where("is_admin = ?", true).Count(&count)
	return count, result.Error
}

func CreateUser(user *models.User) error {
	return portfolio.Data.GetDB().Create(user).Error
}

func HandleGetUsers(c *gin.Context) {
	users := GetUsers()
	c.JSON(http.StatusOK, users)
}

func HandleCreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only a logged-in admin may create another admin
	if input.IsAdmin {
		if u := currentUser(c); u == nil || !u.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required to create an admin user"})
			return
		}
	}

	email := strings.TrimSpace(input.Email)
	if _, err := GetUserByEmail(email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user := &models.User{
		Name:     strings.TrimSpace(input.Name),
		Email:    email,
		Password: portfolio.Encrypt(input.Password),
		IsAdmin:  input.IsAdmin,
	}
	if err := CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func HandleAuth(c *gin.Context) {
	var loginData struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GetUserByName(loginData.Name)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !portfolio.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if err := setAuthCookie(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func HandleValidate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func setAuthCookie(c *gin.Context, user *models.User) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	tokenString, err := token.SignedString([]byte(portfolio.SecretToken))
	if err != nil {
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)
	return nil
}
