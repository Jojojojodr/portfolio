package controllers

import (
	"errors"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"

	"gorm.io/gorm"
)

var ErrLastAdmin = errors.New("cannot remove the last admin user")

func DeleteUser(user *models.User) error {
	if user.IsAdmin {
		adminCount, err := CountAdminUsers()
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return ErrLastAdmin
		}
	}

	return portfolio.Data.GetDB().Transaction(func(tx *gorm.DB) error {
		var postIDs []uint
		if err := tx.Model(&models.BlogPost{}).Where("user_id = ?", user.ID).Pluck("id", &postIDs).Error; err != nil {
			return err
		}
		var commentIDs []uint
		if err := tx.Model(&models.BlogComment{}).Where("user_id = ?", user.ID).Pluck("id", &commentIDs).Error; err != nil {
			return err
		}
		if len(postIDs) > 0 {
			var ids []uint
			if err := tx.Model(&models.BlogComment{}).Where("blog_post_id IN ?", postIDs).Pluck("id", &ids).Error; err != nil {
				return err
			}
			commentIDs = append(commentIDs, ids...)
		}
		if err := deleteComments(tx, commentIDs); err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.CommentLike{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.PostLike{}).Error; err != nil {
			return err
		}
		if len(postIDs) > 0 {
			if err := tx.Where("blog_post_id IN ?", postIDs).Delete(&models.PostLike{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", postIDs).Delete(&models.BlogPost{}).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&models.User{}, user.ID).Error
	})
}

func DeleteBlogPost(postID uint) error {
	return portfolio.Data.GetDB().Transaction(func(tx *gorm.DB) error {
		var commentIDs []uint
		if err := tx.Model(&models.BlogComment{}).Where("blog_post_id = ?", postID).Pluck("id", &commentIDs).Error; err != nil {
			return err
		}
		if err := deleteComments(tx, commentIDs); err != nil {
			return err
		}
		if err := tx.Where("blog_post_id = ?", postID).Delete(&models.PostLike{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.BlogPost{}, postID).Error
	})
}

func DeleteComment(commentID uint) error {
	return portfolio.Data.GetDB().Transaction(func(tx *gorm.DB) error {
		return deleteComments(tx, []uint{commentID})
	})
}

func deleteComments(tx *gorm.DB, commentIDs []uint) error {
	if len(commentIDs) == 0 {
		return nil
	}
	if err := tx.Where("comment_id IN ?", commentIDs).Delete(&models.CommentLike{}).Error; err != nil {
		return err
	}
	return tx.Where("id IN ?", commentIDs).Delete(&models.BlogComment{}).Error
}
