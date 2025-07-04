package views

import (
	"net/http"
	"strconv"

	"github.com/Jojojojodr/portfolio/frontend/blog"
	"github.com/Jojojojodr/portfolio/frontend/components"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	"github.com/gin-gonic/gin"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func HandleBlogPostsPage(c *gin.Context) {
    renderTempl(c, 200, blog.BlogPosts(c))
}

func HandleBlogPostPage(c *gin.Context) {
    idStr := c.Query("id")
	posts, err := models.GetBlogPosts()
	if err != nil {
		c.String(404, "Post not found")
		return
	}

	var post *models.BlogPost = nil

    if idStr != "" {
        id, err := strconv.Atoi(idStr)
        if err == nil {
            post, _ = models.GetBlogPostByID(db.DataBase, uint(id))
			htmlContent := markdown.ToHTML([]byte(post.Content), nil, html.NewRenderer(html.RendererOptions{}))
			post.Content = string(htmlContent)
        }
    }

    renderTempl(c, 200, blog.BlogPostPage(c, posts, post))
}

func HandleBlogPostsHTMX(c *gin.Context) {
    posts, err := models.GetPublishedBlogPosts(db.DataBase)
    if err != nil {
        c.String(500, "Error loading posts")
        return
    }

    c.Writer.WriteHeader(200)
    if len(posts) == 0 {
        components.BlogPostsPartial([]models.BlogPost{}).Render(c.Request.Context(), c.Writer)
    } else {
        components.BlogPostsPartial(posts).Render(c.Request.Context(), c.Writer)
    }
}

func HandleBlogPostHTMX(c *gin.Context) {
    idStr := c.Query("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.Writer.WriteHeader(http.StatusBadRequest)
        c.Writer.Write([]byte("Invalid post ID"))
        return
    }

    post, err := models.GetBlogPostByID(db.DataBase, uint(id))
    if err != nil || post == nil {
        c.Writer.WriteHeader(http.StatusNotFound)
        c.Writer.Write([]byte("Post not found"))
        return
    }

	htmlContent := markdown.ToHTML([]byte(post.Content), nil, html.NewRenderer(html.RendererOptions{}))
	post.Content = string(htmlContent)

    components.BlogPostContent(post).Render(c.Request.Context(), c.Writer)
}