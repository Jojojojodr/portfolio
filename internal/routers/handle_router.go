package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers/views"

	"github.com/gin-gonic/gin"
)

func HandleRouter(svr *gin.Engine) *gin.Engine {
	handle := svr.Group("/handle")
	
	blog := handle.Group("/blog")
	blog.GET("/posts", views.HandleBlogPostsHTMX)

	auth := svr.Group("/auth")
	auth.POST("/login", views.Login)
	auth.POST("/logout", views.Logout)

	return svr
}