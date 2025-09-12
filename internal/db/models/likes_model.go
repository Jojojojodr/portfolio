package models

import (
	"time"

	"github.com/Jojojojodr/portfolio/internal/db"
)

type PostLike struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    UserID     uint      `json:"user_id" gorm:"not null"`
    BlogPostID uint      `json:"blog_post_id" gorm:"not null"`
    CreatedAt  time.Time `json:"created_at"`
    
    // Relationships
    User     User     `json:"-" gorm:"foreignKey:UserID"`
    BlogPost BlogPost `json:"-" gorm:"foreignKey:BlogPostID"`
}

type CommentLike struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    CommentID uint      `json:"comment_id" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    
    // Relationships
    User    User        `json:"-" gorm:"foreignKey:UserID"`
    Comment BlogComment `json:"-" gorm:"foreignKey:CommentID"`
}

// PostLike methods
func CreatePostLike(userID, postID uint) error {
    like := PostLike{
        UserID:     userID,
        BlogPostID: postID,
    }
    return db.DataBase.Create(&like).Error
}

func DeletePostLike(userID, postID uint) error {
    return db.DataBase.Where("user_id = ? AND blog_post_id = ?", userID, postID).Delete(&PostLike{}).Error
}

func IsPostLikedByUser(userID, postID uint) bool {
    var count int64
    db.DataBase.Model(&PostLike{}).Where("user_id = ? AND blog_post_id = ?", userID, postID).Count(&count)
    return count > 0
}

func GetPostLikeCount(postID uint) int64 {
    var count int64
    db.DataBase.Model(&PostLike{}).Where("blog_post_id = ?", postID).Count(&count)
    return count
}

// CommentLike methods
func CreateCommentLike(userID, commentID uint) error {
    like := CommentLike{
        UserID:    userID,
        CommentID: commentID,
    }
    return db.DataBase.Create(&like).Error
}

func DeleteCommentLike(userID, commentID uint) error {
    return db.DataBase.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&CommentLike{}).Error
}

func IsCommentLikedByUser(userID, commentID uint) bool {
    var count int64
    db.DataBase.Model(&CommentLike{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count)
    return count > 0
}

func GetCommentLikeCount(commentID uint) int64 {
    var count int64
    db.DataBase.Model(&CommentLike{}).Where("comment_id = ?", commentID).Count(&count)
    return count
}