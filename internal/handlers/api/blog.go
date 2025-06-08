package api

import (
    "net/http"

    "github.com/Jojojojodr/portfolio/internal/db"	
    "github.com/Jojojojodr/portfolio/internal/db/models"

    "github.com/gin-gonic/gin"
)

func GetPublishedBlogs(c *gin.Context) {
    posts, err := models.GetPublishedBlogPosts(db.DataBase)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blog posts"})
        return
    }
    c.JSON(http.StatusOK, posts)
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