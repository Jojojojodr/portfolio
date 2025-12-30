package routers

import (
	"github.com/Jojojojodr/portfolio/internal/handlers/views"
	"github.com/Jojojojodr/portfolio/internal/middleware"

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

	admin := handle.Group("/admin").Use(middleware.AuthRequired(), middleware.AdminRequired())
	admin.POST("/post/create", views.HandleAdminCreateBlogPost)
	admin.POST("/post/edit", views.HandleAdminEditBlogPost)
	admin.POST("/preview-title", views.HandleAdminPreviewTitle)
	admin.POST("/preview-markdown", views.HandleAdminPreviewMarkdown)

	return svr
}