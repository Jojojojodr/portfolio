package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers"
	"github.com/Jojojojodr/portfolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func V1Router(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the v1 API
	v1 := svr.Group("/v1")
	v1.GET("/health", handlers.HandleHealth)
	v1.GET("/users", handlers.HandleGetUsers)
	v1.POST("/users", handlers.HandleCreateUser)
	v1.POST("/login", handlers.HandleAuth)
	v1.GET("/validate", middleware.AuthMiddleware ,handlers.HandleValidate)

	return svr
}