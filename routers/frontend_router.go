package routers

import (
	"github.com/Jojojojodr/portfolio/middleware"
	"github.com/Jojojojodr/portfolio/views"

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

	blog := svr.Group("/blog")
	blog.GET("/", views.HandleBlogPostsPage)
	blog.GET("/post", views.HandleBlogPostPage)

	svr.NoRoute(views.HandleNotFoundPage)

	return svr
}

func AdminFrontendRouter(svr *gin.Engine) *gin.Engine {
	svr.Static("/static", "./static")
	svr.Use(middleware.LoginMiddleware)

	svr.GET("/login", views.HandleLoginPage)

	admin := svr.Group("/").Use(middleware.AuthRequired(), middleware.AdminRequired())
	admin.GET("/", views.HandleAdminDashboardPage)
	admin.GET("/users", views.HandleAdminUsersPage)
	admin.GET("/user/create", views.HandleAdminCreateUserPage)
	admin.GET("/posts", views.HandleAdminPostsPage)
	admin.GET("/post/create", views.HandleAdminCreateBlogPostPage)
	admin.GET("/post/edit", views.HandleAdminEditBlogPostPage)
	admin.GET("/portfolio", views.HandleAdminPortfolioPage)
	admin.GET("/portfolio/skills", views.HandleAdminPortfolioSkillsPage)
	admin.GET("/portfolio/education", views.HandleAdminPortfolioEducationPage)
	admin.GET("/portfolio/hobbies", views.HandleAdminPortfolioHobbiesPage)
	admin.GET("/portfolio/images", views.HandleAdminPortfolioImagesPage)
	admin.GET("/profile/:id", middleware.AuthRequired(), middleware.AdminRequired(), views.ProfileHandler)

	svr.NoRoute(views.HandleNotFoundPage)

	return svr
}
