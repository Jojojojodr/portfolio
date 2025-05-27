package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers"
	"github.com/Jojojojodr/portfolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func FrontendRouter(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the frontend
	svr.Use(middleware.LoginMiddleware)
	svr.GET("/", handlers.HandleHomePage)
	svr.GET("/login", handlers.HandleLoginPage)

	auth := svr.Group("/auth")
	auth.POST("/login", handlers.Login)
	auth.POST("/logout", handlers.Logout)

	admin := svr.Group("/admin").Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.GET("/dashboard", handlers.HandleDashboardPage)

	return svr
}