package views

import (
	"github.com/Jojojojodr/portfolio/controllers"
	"github.com/Jojojojodr/portfolio/frontend"
	"github.com/Jojojojodr/portfolio/frontend/auth"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func HandleHomePage(c *gin.Context) {
	portfolio, err := controllers.GetPortfolio()
	if err != nil {
		portfolio = nil
	}
	renderTempl(c, 200, frontend.Index(c, portfolio))
}

func HandleLoginPage(c *gin.Context) {
	renderTempl(c, 200, auth.Login(c))
}

func HandleNotFoundPage(c *gin.Context) {
	renderTempl(c, 404, frontend.NotFound(c))
}

func renderTempl(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}
