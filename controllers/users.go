package controllers

import (
	"log"
	"net/http"
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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password before saving
	hashedPassword := portfolio.Encrypt(user.Password)
	user.Password = hashedPassword

	result := portfolio.Data.GetDB().Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func HandleAuth(c *gin.Context) {
	var loginData struct {
		Name     string `json:"name"`
		Password string `json:"password"`
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	secretToken := portfolio.SecretToken
	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func HandleValidate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"user": user})
}
