package views

import (
	"net/http"
	"time"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
	name := c.PostForm("name")
    password := c.PostForm("password")

	if name == "" || password == "" {
        c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "Username and password required"})
        return
    }

    user, err := models.GetUserByName(name)
    if err != nil || user == nil {
        c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "Invalid username or password"})
        return
    }

    if !internal.CheckPasswordHash(password, user.Password) {
        c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "Invalid username or password"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    })

    secretToken := internal.Env("SECRET_TOKEN")
    tokenString, err := token.SignedString([]byte(secretToken))
    if err != nil {
        c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"error": "Could not create token"})
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

    c.Redirect(http.StatusSeeOther, "/")

}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "/", "", false, true) // Clear the cookie

	c.Redirect(http.StatusSeeOther, "/")
}