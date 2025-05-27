package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretToken := internal.Env("SECRET_TOKEN")
		return []byte(secretToken), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims["sub"] == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user *models.User
		user, err := models.GetUserByID(db.DataBase, uint(claims["sub"].(float64)), user)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	
		c.Set("user", user)
	
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func LoginMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.Set("user", nil) // No user if no token
		c.Next()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretToken := internal.Env("SECRET_TOKEN")
		return []byte(secretToken), nil
	})

	if err != nil || !token.Valid {
		c.Set("user", nil) // No user if token is invalid
		c.Next()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		c.Set("user", nil) // No user if claims are invalid
		c.Next()
		return
	}

	var user *models.User
	user, err = models.GetUserByID(db.DataBase, uint(claims["sub"].(float64)), user)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
		c.Set("user", nil) // No user if fetching fails
		c.Next()
		return
	}

	if user == nil {
		c.Set("user", nil) // No user if user is nil
		c.Next()
		return
	}

	c.Set("user", user) // Set user if everything is valid
	c.Next()
}

func IsAuthenticated(c *gin.Context) bool {
	user, exists := c.Get("user")
	return exists && user != nil
}

func AdminMiddleware(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if u, ok := user.(*models.User); ok && u.IsAdmin {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}