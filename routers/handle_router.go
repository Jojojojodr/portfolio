package routers

import (
	"github.com/Jojojojodr/portfolio/views"
	"github.com/Jojojojodr/portfolio/middleware"

	"github.com/gin-gonic/gin"
)

func HandleRouter(svr *gin.Engine) *gin.Engine {
	handle := svr.Group("/handle")
	handle.POST("/register", views.HandleRegisterPost)
	handle.POST("/profile/update", middleware.AuthRequired(), views.UpdateProfileHandler)

	blog := handle.Group("/blog")
	blog.GET("/posts", views.HandleBlogPostsHTMX)
	blog.GET("/post", views.HandleBlogPostHTMX)
	blog.GET("comments", views.HandleGetBlogComments)
	blog.POST("comments/add", views.HandleAddBlogComment)

	like := handle.Group("/like", middleware.AuthRequired())
	like.POST("/post/:postId", views.TogglePostLike)
	like.POST("/comment/:commentId", views.ToggleCommentLike)

	auth := handle.Group("/auth")
	auth.POST("/login", views.Login)
	auth.POST("/logout", views.Logout)

	return svr
}

func AdminHandleRouter(svr *gin.Engine) *gin.Engine {
	handle := svr.Group("/handle")
	handle.Use(middleware.LoginMiddleware)
	handle.POST("/profile/update", middleware.AuthRequired(), views.UpdateProfileHandler)

	auth := handle.Group("/auth")
	auth.POST("/login", views.Login)
	auth.POST("/logout", views.Logout)

	admin := handle.Group("/admin").Use(middleware.AuthRequired(), middleware.AdminRequired())
	admin.DELETE("/user/delete", views.HandleAdminDeleteUser)
	admin.DELETE("/post/delete", views.HandleAdminDeleteBlogPost)
	admin.POST("/user/create", views.HandleAdminCreateUser)
	admin.POST("/post/create", views.HandleAdminCreateBlogPost)
	admin.POST("/post/edit", views.HandleAdminEditBlogPost)
	admin.POST("/preview-title", views.HandleAdminPreviewTitle)
	admin.POST("/preview-markdown", views.HandleAdminPreviewMarkdown)

	portfolio := handle.Group("/admin/portfolio").Use(middleware.AuthRequired(), middleware.AdminRequired())
	portfolio.POST("/profile", views.HandleAdminSavePortfolioProfile)
	portfolio.POST("/about", views.HandleAdminSavePortfolioAbout)
	portfolio.POST("/skill/save", views.HandleAdminSaveSkill)
	portfolio.POST("/skill/delete", views.HandleAdminDeleteSkill)
	portfolio.POST("/soft-skill/save", views.HandleAdminSaveSoftSkill)
	portfolio.POST("/soft-skill/delete", views.HandleAdminDeleteSoftSkill)
	portfolio.POST("/language/save", views.HandleAdminSaveLanguage)
	portfolio.POST("/language/delete", views.HandleAdminDeleteLanguage)
	portfolio.POST("/education/save", views.HandleAdminSaveEducation)
	portfolio.POST("/education/delete", views.HandleAdminDeleteEducation)
	portfolio.POST("/education/task/save", views.HandleAdminSaveEducationTask)
	portfolio.POST("/education/task/delete", views.HandleAdminDeleteEducationTask)
	portfolio.POST("/hobby/save", views.HandleAdminSaveHobby)
	portfolio.POST("/hobby/delete", views.HandleAdminDeleteHobby)

	return svr
}
