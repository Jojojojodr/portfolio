package models

import (
	"time"

	"github.com/Jojojojodr/portfolio/internal/db"
	"gorm.io/gorm"
)

type BlogPost struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      uint      `json:"user_id"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetPublishedBlogPosts(db *gorm.DB) ([]BlogPost, error) {
	var posts []BlogPost
	err := db.Where("is_published = ?", 1).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetBlogPosts() ([]BlogPost, error) {
	var posts []BlogPost
	err := db.DataBase.Where("is_published = ?", 1).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}