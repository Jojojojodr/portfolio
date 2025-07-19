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
	User		User      `gorm:"foreignKey:UserID"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

type BlogComment struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    Comment    string    `json:"comment"`
    UserID     uint      `json:"user_id"`
    User       User      `gorm:"foreignKey:UserID"`
    BlogPostID uint      `json:"blog_post_id"`
    BlogPost   BlogPost  `gorm:"foreignKey:BlogPostID"`
    CreatedAt  time.Time `json:"created_at"`
}

func (bp *BlogPost) GetLikeCount() int64 {
    return GetPostLikeCount(bp.ID)
}

func (bp *BlogPost) IsLikedByUser(userID uint) bool {
    return IsPostLikedByUser(userID, bp.ID)
}

func (bc *BlogComment) GetLikeCount() int64 {
    return GetCommentLikeCount(bc.ID)
}

func (bc *BlogComment) IsLikedByUser(userID uint) bool {
    return IsCommentLikedByUser(userID, bc.ID)
}

func GetPublishedBlogPosts(db *gorm.DB) ([]BlogPost, error) {
	var posts []BlogPost
	err := db.Preload("User").Where("is_published = ?", 1).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetBlogPosts() ([]BlogPost, error) {
	var posts []BlogPost
	err := db.DataBase.Preload("User").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetBlogPostByID(db *gorm.DB, id uint) (*BlogPost, error) {
    var post BlogPost
    err := db.Preload("User").First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}