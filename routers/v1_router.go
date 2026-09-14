package routers

import (
	"github.com/Jojojojodr/portfolio/controllers"
	"github.com/Jojojojodr/portfolio/middleware"

	"github.com/gin-gonic/gin"
)

func V1Router(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the v1 API
	v1 := svr.Group("/v1")
	v1.GET("/health", controllers.HandleHealth)
	v1.GET("/users", controllers.HandleGetUsers)
	v1.GET("/validate", middleware.AuthMiddleware, controllers.HandleValidate)

	v1.POST("/users", controllers.HandleCreateUser)
	v1.POST("/login", controllers.HandleAuth)
	
	blog := v1.Group("/blog")
	blog.GET("/posts", controllers.PublishedBlogs)
	blog.GET("/post", controllers.GetBlogPostByID)
	
	blog.POST("/create-post", middleware.AuthMiddleware, middleware.AdminMiddleware, controllers.CreateBlogPost)

	return svr
}