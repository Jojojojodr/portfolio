package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Jojojojodr/portfolio/frontend/admin"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"github.com/gin-gonic/gin"
)

func HandleAdminDashboardPage(c *gin.Context) {
	users := models.GetUsers()

	renderTempl(c, 200, admin.Dashboard(c, users))
}

func HandleAdminUsersPage(c *gin.Context) {
    users := models.GetUsers()
    renderTempl(c, 200, admin.AdminUsersPage(c, users))
}

func HandleAdminPostsPage(c *gin.Context) {
    posts, err := models.GetBlogPosts()
    if err != nil {
        c.String(500, "Failed to load posts")
        return
    }
    renderTempl(c, 200, admin.AdminPostsPage(c, posts))
}

func HandleAdminCreateBlogPostPage(c *gin.Context) {
    admin.BlogCreatePage(c, "").Render(c.Request.Context(), c.Writer)
}

func HandleAdminCreateBlogPost(c *gin.Context) {
    title := c.PostForm("title")
    content := c.PostForm("content")
    isPublished := c.PostForm("is_published") == "1"

    userAny, _ := c.Get("user")
    user, ok := userAny.(*models.User)
    if !ok || user == nil {
        admin.BlogCreatePage(c, "You must be logged in.").Render(c.Request.Context(), c.Writer)
        return
    }

    post := models.BlogPost{
        Title:       	title,
        Content:     	content,
        UserID:      	user.ID,
        IsPublished: 	isPublished,
		CreatedAt: 		time.Now(),
    }
    if err := db.DataBase.Create(&post).Error; err != nil {
        admin.BlogCreatePage(c, "Failed to create post.").Render(c.Request.Context(), c.Writer)
        return
    }
    c.Redirect(302, "/admin/posts")
}

func HandleAdminEditBlogPostPage(c *gin.Context) {
    idStr := c.Query("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid post ID")
        return
    }
    post, err := models.GetBlogPostByID(uint(id))
    if err != nil || post == nil {
        c.String(http.StatusNotFound, "Post not found")
        return
    }
    admin.BlogEditPage(c, post, "").Render(c.Request.Context(), c.Writer)
}

func HandleAdminEditBlogPost(c *gin.Context) {
    idStr := c.Query("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid post ID")
        return
    }
    post, err := models.GetBlogPostByID(uint(id))
    if err != nil || post == nil {
        c.String(http.StatusNotFound, "Post not found")
        return
    }

    title := c.PostForm("title")
    content := c.PostForm("content")
    isPublished := c.PostForm("is_published") == "1"

    if title == "" || content == "" {
        admin.BlogEditPage(c, post, "Title and content are required.").Render(c.Request.Context(), c.Writer)
        return
    }

    post.Title = title
    post.Content = content
    post.IsPublished = isPublished

    if err := db.DataBase.Save(post).Error; err != nil {
        admin.BlogEditPage(c, post, "Failed to update post.").Render(c.Request.Context(), c.Writer)
        return
    }

    c.Redirect(http.StatusSeeOther, "/admin/posts")
}