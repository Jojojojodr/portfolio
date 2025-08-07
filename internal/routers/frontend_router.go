package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers/views"
	"github.com/Jojojojodr/portfolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func FrontendRouter(svr *gin.Engine) *gin.Engine {
	// Set up the routes for the frontend
	svr.Static("/static", "./static")

	svr.Use(middleware.LoginMiddleware)
	svr.GET("/", views.HandleHomePage)
	svr.GET("/login", views.HandleLoginPage)
	svr.GET("/register", views.HandleRegisterPage)
	svr.GET("/profile", middleware.AuthRequired(), views.ProfileHandler)
	svr.GET("/profile/:id", middleware.AuthRequired(), middleware.AdminRequired(), views.ProfileHandler)

	blog := svr.Group("/blog")
	blog.GET("/", views.HandleBlogPostsPage)
	blog.GET("/post", views.HandleBlogPostPage)

	admin := svr.Group("/admin").Use(middleware.AuthRequired(), middleware.AdminRequired())
	admin.GET("/dashboard", views.HandleAdminDashboardPage)
	admin.GET("/users", views.HandleAdminUsersPage)
	admin.GET("/posts", views.HandleAdminPostsPage)
	admin.GET("/post/create", views.HandleAdminCreateBlogPostPage)
	admin.GET("/post/edit", views.HandleAdminEditBlogPostPage)
	admin.POST("/post/create", views.HandleAdminCreateBlogPost)
	admin.POST("/post/edit", views.HandleAdminEditBlogPost)

	svr.NoRoute(views.HandleNotFoundPage)

	return svr
}