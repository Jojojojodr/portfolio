package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers/views"
	"github.com/Jojojojodr/portfolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func FrontendRouter(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the frontend
	svr.Use(middleware.LoginMiddleware)
	svr.GET("/", views.HandleHomePage)
	svr.GET("/login", views.HandleLoginPage)

	auth := svr.Group("/auth")
	auth.POST("/login", views.Login)
	auth.POST("/logout", views.Logout)

	admin := svr.Group("/admin").Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.GET("/dashboard", views.HandleDashboardPage)

	return svr
}