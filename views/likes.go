package views

import (
    "net/http"
    "strconv"

	"github.com/Jojojojodr/portfolio/controllers"
    "github.com/Jojojojodr/portfolio/database/models"
    "github.com/gin-gonic/gin"
)

func TogglePostLike(c *gin.Context) {
    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    user := userInterface.(*models.User)
    postIDStr := c.Param("postId")
    postID, err := strconv.ParseUint(postIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Check if already liked
    if controllers.IsPostLikedByUser(user.ID, uint(postID)) {
        // Unlike
        err = controllers.DeletePostLike(user.ID, uint(postID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
            return
        }
    } else {
        // Like
        err = controllers.CreatePostLike(user.ID, uint(postID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
            return
        }
    }

    // Return updated data
    likeCount := controllers.GetPostLikeCount(uint(postID))
    isLiked := controllers.IsPostLikedByUser(user.ID, uint(postID))

    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "like_count": likeCount,
        "is_liked":   isLiked,
    })
}

func ToggleCommentLike(c *gin.Context) {
    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    user := userInterface.(*models.User)
    commentIDStr := c.Param("commentId")
    commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // Check if already liked
    if controllers.IsCommentLikedByUser(user.ID, uint(commentID)) {
        // Unlike
        err = controllers.DeleteCommentLike(user.ID, uint(commentID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike comment"})
            return
        }
    } else {
        // Like
        err = controllers.CreateCommentLike(user.ID, uint(commentID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like comment"})
            return
        }
    }

    // Return updated data
    likeCount := controllers.GetCommentLikeCount(uint(commentID))
    isLiked := controllers.IsCommentLikedByUser(user.ID, uint(commentID))

    c.JSON(http.StatusOK, gin.H{
        "success":    true,
        "like_count": likeCount,
        "is_liked":   isLiked,
    })
}