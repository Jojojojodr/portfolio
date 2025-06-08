package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers/api"
	"github.com/Jojojojodr/portfolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func V1Router(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the v1 API
	v1 := svr.Group("/v1")
	v1.GET("/health", api.HandleHealth)
	v1.GET("/users", api.HandleGetUsers)
	v1.POST("/users", api.HandleCreateUser)
	v1.POST("/login", api.HandleAuth)
	v1.GET("/validate", middleware.AuthMiddleware ,api.HandleValidate) 

	
	blog := v1.Group("/blog")
	blog.GET("/posts", api.GetPublishedBlogs)
	blog.POST("/create-post", middleware.AuthMiddleware, middleware.AdminMiddleware, api.CreateBlogPost)

	return svr
}