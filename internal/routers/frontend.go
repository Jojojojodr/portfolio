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

	blog := svr.Group("/blog")
	blog.GET("/", views.HandleBlogPostsPage)
	blog.GET("/post/:id", views.HandleBlogPostPage)

	admin := svr.Group("/admin").Use(middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.GET("/dashboard", views.HandleDashboardPage)
	admin.GET("/blog/create", views.HandleCreateBlogPostPage)
	admin.POST("/blog/create", views.HandleCreateBlogPost)

	svr.NoRoute(views.HandleNotFoundPage)

	return svr
}