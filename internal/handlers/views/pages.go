package views

import (
	"github.com/Jojojojodr/portfolio/frontend"
	"github.com/Jojojojodr/portfolio/frontend/auth"
	"github.com/Jojojojodr/portfolio/frontend/admin"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func HandleHomePage(c *gin.Context) {
    renderTempl(c, 200, frontend.Index(c))
}

func HandleLoginPage(c *gin.Context) {
	renderTempl(c, 200, auth.Login(c))
}

func HandleDashboardPage(c *gin.Context) {
	users := models.GetUsers()

	renderTempl(c, 200, admin.Dashboard(c, users))
}

func renderTempl(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}