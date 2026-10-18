package views

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/controllers"
	"github.com/Jojojojodr/portfolio/frontend/admin"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func HandleAdminDashboardPage(c *gin.Context) {
	renderTempl(c, 200, admin.Dashboard(c))
}

func HandleAdminUsersPage(c *gin.Context) {
    users := controllers.GetUsers()
    renderTempl(c, 200, admin.AdminUsersPage(c, users))
}

func HandleAdminCreateUserPage(c *gin.Context) {
	errors := make(map[string]string)
	renderTempl(c, 200, admin.UserCreatePage(c, errors))
}

func HandleAdminCreateUser(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	isAdmin := c.PostForm("is_admin") == "on"

	errors := make(map[string]string)
	if name == "" {
		errors["name"] = "Name is required"
	}
	if email == "" {
		errors["email"] = "Email is required"
	} else if !strings.Contains(email, "@") {
		errors["email"] = "Invalid email"
	}
	if password == "" {
		errors["password"] = "Password is required"
	} else if len(password) < 6 {
		errors["password"] = "Password must be at least 6 characters"
	}
	if confirmPassword == "" {
		errors["confirm_password"] = "Password confirmation is required"
	} else if password != confirmPassword {
		errors["confirm_password"] = "Passwords do not match"
	}
	if _, err := controllers.GetUserByEmail(email); err == nil {
		errors["email"] = "Email already registered"
	}

	if len(errors) > 0 {
		admin.UserCreatePage(c, errors).Render(c.Request.Context(), c.Writer)
		return
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: portfolio.Encrypt(password),
		IsAdmin:  isAdmin,
	}
	if err := controllers.CreateUser(user); err != nil {
		errors["general"] = "Failed to create user"
		admin.UserCreatePage(c, errors).Render(c.Request.Context(), c.Writer)
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")
}

func HandleAdminDeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := controllers.GetUserByID(uint(id))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := controllers.DeleteUser(user); err != nil {
		if errors.Is(err, controllers.ErrLastAdmin) {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete the last admin user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func HandleAdminDeleteBlogPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	post, err := controllers.BlogPostByID(uint(id))
	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := controllers.DeleteBlogPost(post.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func HandleAdminPostsPage(c *gin.Context) {
	posts, err := controllers.BlogPosts()
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
		Title:       title,
		Content:     content,
		UserID:      user.ID,
		IsPublished: isPublished,
		CreatedAt:   time.Now(),
	}
	if err := portfolio.Data.GetDB().Create(&post).Error; err != nil {
		admin.BlogCreatePage(c, "Failed to create post.").Render(c.Request.Context(), c.Writer)
		return
	}
	c.Redirect(302, "/posts")
}

func HandleAdminEditBlogPostPage(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid post ID")
		return
	}
	post, err := controllers.BlogPostByID(uint(id))
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
	post, err := controllers.BlogPostByID(uint(id)) 
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

	if err := portfolio.Data.GetDB().Save(post).Error; err != nil {
		admin.BlogEditPage(c, post, "Failed to update post.").Render(c.Request.Context(), c.Writer)
		return
	}

	c.Redirect(http.StatusSeeOther, "/posts")
}

func HandleAdminPreviewTitle(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		title = "Blog Post Title"
	}
	c.String(200, title)
}

func HandleAdminPreviewMarkdown(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.String(200, `<p class="text-gray-400">Your markdown content will appear here as you type...</p>`)
		return
	}

	// Parse markdown to HTML using the same method as blog posts
	htmlContent := markdown.ToHTML([]byte(content), nil, html.NewRenderer(html.RendererOptions{}))

	// Wrap the HTML content in the same structure as the MarkdownRender component
	wrappedHTML := `<div class="markdown">` + string(htmlContent) + `</div>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, wrappedHTML)
}