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
	v1.GET("/validate", middleware.AuthMiddleware, controllers.HandleValidate)

	v1.POST("/register", controllers.HandleCreateUser)
	v1.POST("/login", controllers.HandleAuth)

	blog := v1.Group("/blog")
	blog.GET("/posts", controllers.PublishedBlogs)
	blog.GET("/post", controllers.GetBlogPostByID)

	portfolio := v1.Group("/portfolio")
	portfolio.GET("", controllers.GetPortfolioAPI)
	portfolio.GET("/skills", controllers.GetSkillsAPI)
	portfolio.GET("/soft-skills", controllers.GetSoftSkillsAPI)
	portfolio.GET("/languages", controllers.GetLanguagesAPI)
	portfolio.GET("/education", controllers.GetEducationAPI)
	portfolio.GET("/hobbies", controllers.GetHobbiesAPI)

	return svr
}

func AdminV1Router(svr *gin.Engine) *gin.Engine {
	// Set up the JSON API routes for the admin server
	v1 := svr.Group("/v1")
	v1.GET("/health", controllers.HandleHealth)
	v1.POST("/login", controllers.HandleAdminAuth)

	admin := v1.Group("", middleware.AuthMiddleware, middleware.AdminMiddleware)
	admin.POST("/logout", controllers.HandleAdminLogout)
	admin.GET("/validate", controllers.HandleValidate)

	admin.GET("/users", controllers.HandleGetUsers)
	admin.GET("/user", controllers.HandleAdminGetUser)
	admin.POST("/user/create", controllers.HandleAdminCreateUser)
	admin.PATCH("/user/edit", controllers.HandleAdminEditUser)
	admin.DELETE("/user/delete", controllers.HandleAdminDeleteUser)

	admin.GET("/posts", controllers.HandleAdminGetPosts)
	admin.GET("/post", controllers.HandleAdminGetPost)
	admin.POST("/post/create", controllers.HandleAdminCreatePost)
	admin.PATCH("/post/edit", controllers.HandleAdminEditPost)
	admin.DELETE("/post/delete", controllers.HandleAdminDeletePost)

	admin.GET("/comments", controllers.HandleAdminGetComments)
	admin.DELETE("/comment/delete", controllers.HandleAdminDeleteComment)

	portfolio := admin.Group("/portfolio")
	portfolio.PATCH("", controllers.HandleAdminUpdatePortfolio)
	portfolio.PUT("/about", controllers.HandleAdminUpdateAbout)
	portfolio.PUT("/skills", controllers.HandleAdminReplaceSkills)
	portfolio.PUT("/soft-skills", controllers.HandleAdminReplaceSoftSkills)
	portfolio.PUT("/languages", controllers.HandleAdminReplaceLanguages)
	portfolio.PUT("/education", controllers.HandleAdminReplaceEducation)
	portfolio.PUT("/hobbies", controllers.HandleAdminReplaceHobbies)

	return svr
}
