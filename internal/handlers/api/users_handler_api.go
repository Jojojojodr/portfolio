package api

import (
	"net/http"
	"time"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func HandleGetUsers(c *gin.Context) {
    users := models.GetUsers()
	c.JSON(http.StatusOK, users)
}

func HandleCreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash the password before saving
	hashedPassword := internal.Encrypt(user.Password)
	user.Password = hashedPassword

	result := db.DataBase.Create(&user)
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

	user, err := models.GetUserByName(loginData.Name)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !internal.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	
	secretToken := internal.Env("SECRET_TOKEN")
	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func HandleValidate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"user": user})
}