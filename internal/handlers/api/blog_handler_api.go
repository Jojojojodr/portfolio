package api

import (
    "net/http"
    "fmt"

    "github.com/Jojojojodr/portfolio/internal/db"	
    "github.com/Jojojojodr/portfolio/internal/db/models"

    "github.com/gin-gonic/gin"
)

func GetPublishedBlogs(c *gin.Context) {
    posts, err := models.GetPublishedBlogPosts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blog posts"})
        return
    }
    c.JSON(http.StatusOK, posts)
}

func GetBlogPostByID(c *gin.Context) {
    idStr := c.Query("id")
    var idUint uint
    _, err := fmt.Sscanf(idStr, "%d", &idUint)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
        return
    }

    post, err := models.GetBlogPostByID(idUint)
    if err != nil {
        c.JSON(404, gin.H{"error": "Blog post not found"})
        return
    }

	comments, err := models.GetCommentsByPostID(post.ID)
	if err != nil {
		comments = []models.BlogComment{}
	}

    c.JSON(200, gin.H{
		"post": post,
		"comments": comments,
	})
}

func CreateBlogPost(c *gin.Context) {
    var postInput struct {
		Title       string `json:"title" binding:"required"`
    	Content     string `json:"content" binding:"required"`
    	UserID      uint   `json:"user_id" binding:"required"`
    	IsPublished bool   `json:"is_published"`
	}

    if err := c.ShouldBindJSON(&postInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    blogPost := models.BlogPost{
        Title:       postInput.Title,
        Content:     postInput.Content,
        UserID:      postInput.UserID,
        IsPublished: postInput.IsPublished,
    }

    if err := db.DataBase.Create(&blogPost).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post"})
        return
    }

    c.JSON(http.StatusCreated, blogPost)
}