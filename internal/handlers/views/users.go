package views

import (
	"net/http"
	"time"

	"github.com/Jojojojodr/portfolio/frontend/components"
	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
    name := c.PostForm("name")
    password := c.PostForm("password")

    if name == "" || password == "" {
		c.Writer.WriteHeader(400)
		components.LoginResponse("", "Username and password required").Render(c.Request.Context(), c.Writer)
        return
    }

    user, err := models.GetUserByName(name)
    if err != nil || user == nil {
		c.Writer.WriteHeader(401)
		components.LoginResponse("", "Invalid username or password").Render(c.Request.Context(), c.Writer)
        return
    }

    if !internal.CheckPasswordHash(password, user.Password) {
		c.Writer.WriteHeader(401)
		components.LoginResponse("", "Invalid username or password").Render(c.Request.Context(), c.Writer)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    })

    secretToken := internal.Env("SECRET_TOKEN")
    tokenString, err := token.SignedString([]byte(secretToken))
    if err != nil {
		c.Writer.WriteHeader(500)
		components.LoginResponse("", "Could not create token").Render(c.Request.Context(), c.Writer)
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	c.Header("HX-Location", "/")
	c.Writer.WriteHeader(200)
}

func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "/", "", false, true) // Clear the cookie

	c.Redirect(http.StatusSeeOther, "/")
}