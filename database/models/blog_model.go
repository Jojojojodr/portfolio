package models

import (
	"time"
)

type BlogPost struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      uint      `json:"user_id"`
	User		User      `json:"-" gorm:"foreignKey:UserID"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

type BlogComment struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    Comment    string    `json:"comment"`
    UserID     uint      `json:"user_id"`
    User       User      `json:"-" gorm:"foreignKey:UserID"`
    BlogPostID uint      `json:"blog_post_id"`
    BlogPost   BlogPost  `json:"-" gorm:"foreignKey:BlogPostID"`
    CreatedAt  time.Time `json:"created_at"`
}