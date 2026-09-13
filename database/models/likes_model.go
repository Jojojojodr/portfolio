package models

import (
	"time"
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