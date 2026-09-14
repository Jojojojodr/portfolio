package controllers

import (
	"net/http"
    "fmt"

    "github.com/Jojojojodr/portfolio"	
    "github.com/Jojojojodr/portfolio/database/models"

    "github.com/gin-gonic/gin"
)

func PublishedBlogPosts() ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := portfolio.Data.GetDB().Preload("User").Where("is_published = true").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func BlogPosts() ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := portfolio.Data.GetDB().Preload("User").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func BlogPostByID(id uint) (*models.BlogPost, error) {
    var post models.BlogPost
    err := portfolio.Data.GetDB().Preload("User").First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func CommentsByPostID(postID uint) ([]models.BlogComment, error) {
    var comments []models.BlogComment
    err := portfolio.Data.GetDB().Where("blog_post_id = ?", postID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
    return comments, err
}

func PublishedBlogs(c *gin.Context) {
    posts, err := PublishedBlogPosts()
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

    post, err := BlogPostByID(idUint)
    if err != nil {
        c.JSON(404, gin.H{"error": "Blog post not found"})
        return
    }

	comments, err := CommentsByPostID(post.ID)
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

    if err := portfolio.Data.GetDB().Create(&blogPost).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post"})
        return
    }

    c.JSON(http.StatusCreated, blogPost)
}