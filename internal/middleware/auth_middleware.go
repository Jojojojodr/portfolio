package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jojojojodr/portfolio/internal"
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
		user, err = models.GetUserByID(uint(claims["sub"].(float64)))
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
        c.Set("isAuthenticated", false)
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
        c.Set("isAuthenticated", false)
        c.Next()
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.Set("isAuthenticated", false)
        c.Next()
        return
    }

    userID, ok := claims["sub"].(float64)
    if !ok {
        c.Set("isAuthenticated", false)
        c.Next()
        return
    }

    user, err := models.GetUserByID(uint(userID))
    if err != nil || user == nil {
        c.Set("isAuthenticated", false)
        c.Next()
        return
    }

    c.Set("isAuthenticated", true)
    c.Set("user", user)
    c.Set("isAdmin", user.IsAdmin)
    
    c.Next()
}

func IsAuthenticated(c *gin.Context) bool {
	user, exists := c.Get("user")
	return exists && user != nil
}

func IsAdmin(c *gin.Context) bool {
	user, exists := c.Get("user")
	if !exists {
		return false
	}

	if u, ok := user.(*models.User); ok {
		return u.IsAdmin
	}
	return false
}

func GetUser(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}

	if u, ok := user.(*models.User); ok {
		return u
	}
	return nil
}

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAuthenticated, exists := c.Get("isAuthenticated")
        if !exists || !isAuthenticated.(bool) {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        _, userExists := c.Get("user")
        if !userExists {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }
        c.Next()
    }
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

func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAuthenticated, exists := c.Get("isAuthenticated")
        if !exists || !isAuthenticated.(bool) {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        isAdmin, adminExists := c.Get("isAdmin")
        if !adminExists || !isAdmin.(bool) {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        c.Next()
    }
}