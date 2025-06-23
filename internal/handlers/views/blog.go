package views

import (
	"time"

	"github.com/Jojojojodr/portfolio/frontend/admin"
	"github.com/Jojojojodr/portfolio/frontend/blog"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	"github.com/gin-gonic/gin"

	"github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
)

func HandleBlogPostsPage(c *gin.Context) {
    posts, err := models.GetBlogPosts()
	if err != nil {
		renderTempl(c, 200, blog.BlogPosts(c, nil))
	}
    renderTempl(c, 200, blog.BlogPosts(c, posts))
}

func HandleBlogPostPage(c *gin.Context) {
    id := c.Param("id")
	posts, err := models.GetBlogPosts()
	if err != nil {
		c.Redirect(404, "/not-found")
        return
	}
    var post models.BlogPost
    if err := db.DataBase.First(&post, id).Error; err != nil {
        c.Redirect(404, "/not-found")
        return
    }
	htmlContent := markdown.ToHTML([]byte(post.Content), nil, html.NewRenderer(html.RendererOptions{}))
	post.Content = string(htmlContent)
    renderTempl(c, 200, blog.BlogPostPage(c, posts, &post))
}

func HandleCreateBlogPostPage(c *gin.Context) {
    admin.BlogCreatePage(c, "").Render(c.Request.Context(), c.Writer)
}

func HandleCreateBlogPost(c *gin.Context) {
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
    c.Redirect(302, "/blog")
}