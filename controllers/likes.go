package controllers

import (
	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"
)

// PostLike methods
func CreatePostLike(userID, postID uint) error {
    like := models.PostLike{
        UserID:     userID,
        BlogPostID: postID,
    }
    return portfolio.Data.GetDB().Create(&like).Error
}

func DeletePostLike(userID, postID uint) error {
    return portfolio.Data.GetDB().Where("user_id = ? AND blog_post_id = ?", userID, postID).Delete(&models.PostLike{}).Error
}

func IsPostLikedByUser(userID, postID uint) bool {
    var count int64
    portfolio.Data.GetDB().Model(&models.PostLike{}).Where("user_id = ? AND blog_post_id = ?", userID, postID).Count(&count)
    return count > 0
}

func GetPostLikeCount(postID uint) int64 {
    var count int64
    portfolio.Data.GetDB().Model(&models.PostLike{}).Where("blog_post_id = ?", postID).Count(&count)
    return count
}

// CommentLike methods
func CreateCommentLike(userID, commentID uint) error {
    like := models.CommentLike{
        UserID:    userID,
        CommentID: commentID,
    }
    return portfolio.Data.GetDB().Create(&like).Error
}

func DeleteCommentLike(userID, commentID uint) error {
    return portfolio.Data.GetDB().Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{}).Error
}

func IsCommentLikedByUser(userID, commentID uint) bool {
    var count int64
    portfolio.Data.GetDB().Model(&models.CommentLike{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count)
    return count > 0
}

func GetCommentLikeCount(commentID uint) int64 {
    var count int64
    portfolio.Data.GetDB().Model(&models.CommentLike{}).Where("comment_id = ?", commentID).Count(&count)
    return count
}