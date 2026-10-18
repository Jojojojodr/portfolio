package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Jojojojodr/portfolio"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleAdminAuth(c *gin.Context) {
	var loginData struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and password are required"})
		return
	}

	user, err := GetUserByName(loginData.Name)
	if err != nil || user == nil || !portfolio.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin privileges required"})
		return
	}

	if err := setAuthCookie(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

func HandleAdminLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func HandleAdminGetUser(c *gin.Context) {
	id, ok := queryID(c, "user_id")
	if !ok {
		return
	}
	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func HandleAdminCreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		IsAdmin  bool   `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.TrimSpace(input.Email)
	if _, err := GetUserByEmail(email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user := &models.User{
		Name:     strings.TrimSpace(input.Name),
		Email:    email,
		Password: portfolio.Encrypt(input.Password),
		IsAdmin:  input.IsAdmin,
	}
	if err := CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func HandleAdminEditUser(c *gin.Context) {
	id, ok := queryID(c, "user_id")
	if !ok {
		return
	}
	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name     *string `json:"name" binding:"omitempty,min=1"`
		Email    *string `json:"email" binding:"omitempty,email"`
		Password *string `json:"password" binding:"omitempty,min=6"`
		IsAdmin  *bool   `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		user.Name = strings.TrimSpace(*input.Name)
	}
	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		if existing, err := GetUserByEmail(email); err == nil && existing.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		user.Email = email
	}
	if input.Password != nil {
		user.Password = portfolio.Encrypt(*input.Password)
	}
	if input.IsAdmin != nil {
		if user.IsAdmin && !*input.IsAdmin {
			adminCount, err := CountAdminUsers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
				return
			}
			if adminCount <= 1 {
				c.JSON(http.StatusConflict, gin.H{"error": ErrLastAdmin.Error()})
				return
			}
		}
		user.IsAdmin = *input.IsAdmin
	}

	if err := UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func HandleAdminDeleteUser(c *gin.Context) {
	id, ok := queryID(c, "user_id")
	if !ok {
		return
	}
	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := DeleteUser(user); err != nil {
		if errors.Is(err, ErrLastAdmin) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func HandleAdminGetPosts(c *gin.Context) {
	posts, err := BlogPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch blog posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func HandleAdminGetPost(c *gin.Context) {
	id, ok := queryID(c, "post_id")
	if !ok {
		return
	}
	post, err := BlogPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	comments, err := CommentsByPostID(post.ID)
	if err != nil {
		comments = []models.BlogComment{}
	}
	c.JSON(http.StatusOK, gin.H{"post": post, "comments": comments})
}

func HandleAdminCreatePost(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		IsPublished bool   `json:"is_published"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.BlogPost{
		Title:       input.Title,
		Content:     input.Content,
		UserID:      currentUser(c).ID,
		IsPublished: input.IsPublished,
		CreatedAt:   time.Now(),
	}
	if err := portfolio.Data.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func HandleAdminEditPost(c *gin.Context) {
	id, ok := queryID(c, "post_id")
	if !ok {
		return
	}
	post, err := BlogPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input struct {
		Title       *string `json:"title" binding:"omitempty,min=1"`
		Content     *string `json:"content" binding:"omitempty,min=1"`
		IsPublished *bool   `json:"is_published"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.IsPublished != nil {
		post.IsPublished = *input.IsPublished
	}

	if err := portfolio.Data.GetDB().Save(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func HandleAdminDeletePost(c *gin.Context) {
	id, ok := queryID(c, "post_id")
	if !ok {
		return
	}
	if _, err := BlogPostByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := DeleteBlogPost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func HandleAdminGetComments(c *gin.Context) {
	id, ok := queryID(c, "post_id")
	if !ok {
		return
	}
	comments, err := CommentsByPostID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func HandleAdminDeleteComment(c *gin.Context) {
	id, ok := queryID(c, "comment_id")
	if !ok {
		return
	}
	var comment models.BlogComment
	if err := portfolio.Data.GetDB().First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}

	if err := DeleteComment(comment.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

func HandleAdminUpdatePortfolio(c *gin.Context) {
	p, ok := loadPortfolio(c)
	if !ok {
		return
	}

	var input struct {
		Name     *string `json:"name" binding:"omitempty,min=1"`
		Function *string `json:"function"`
		Place    *string `json:"place"`
		Email    *string `json:"email" binding:"omitempty,email"`
		Github   *string `json:"github" binding:"omitempty,url"`
		LinkedIn *string `json:"linkedin" binding:"omitempty,url"`
		Birthday *string `json:"birthday"`
		ImgPath  *string `json:"img_path"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fields := []struct {
		value  *string
		target *string
	}{
		{input.Name, &p.Name},
		{input.Function, &p.Function},
		{input.Place, &p.Place},
		{input.Email, &p.Email},
		{input.Github, &p.Github},
		{input.LinkedIn, &p.LinkedIn},
		{input.Birthday, &p.Birthday},
		{input.ImgPath, &p.ImgPath},
	}
	for _, f := range fields {
		if f.value != nil {
			*f.target = strings.TrimSpace(*f.value)
		}
	}

	savePortfolio(c, p)
}

func HandleAdminUpdateAbout(c *gin.Context) {
	p, ok := loadPortfolio(c)
	if !ok {
		return
	}

	var input struct {
		About string `json:"about" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "about is required"})
		return
	}

	p.About = strings.TrimSpace(input.About)
	savePortfolio(c, p)
}

func HandleAdminReplaceSkills(c *gin.Context) {
	replacePortfolioList(c, "skills", func(s *models.Skill, portfolioID string) ([]string, error) {
		if s.Skill = strings.TrimSpace(s.Skill); s.Skill == "" {
			return nil, errors.New("skill is required")
		}
		s.PortfolioID = portfolioID
		return []string{ensureID(&s.ID, "skill")}, nil
	}, nil)
}

func HandleAdminReplaceSoftSkills(c *gin.Context) {
	replacePortfolioList(c, "soft_skills", func(s *models.SoftSkill, portfolioID string) ([]string, error) {
		if s.Skill = strings.TrimSpace(s.Skill); s.Skill == "" {
			return nil, errors.New("skill is required")
		}
		s.PortfolioID = portfolioID
		return []string{ensureID(&s.ID, "soft-skill")}, nil
	}, nil)
}

func HandleAdminReplaceLanguages(c *gin.Context) {
	replacePortfolioList(c, "languages", func(l *models.Language, portfolioID string) ([]string, error) {
		if l.Language = strings.TrimSpace(l.Language); l.Language == "" {
			return nil, errors.New("language is required")
		}
		l.PortfolioID = portfolioID
		return []string{ensureID(&l.ID, "language")}, nil
	}, nil)
}

func HandleAdminReplaceEducation(c *gin.Context) {
	replacePortfolioList(c, "education", func(e *models.Education, portfolioID string) ([]string, error) {
		if e.Education = strings.TrimSpace(e.Education); e.Education == "" {
			return nil, errors.New("education is required")
		}
		if e.School = strings.TrimSpace(e.School); e.School == "" {
			return nil, errors.New("school is required")
		}
		e.PortfolioID = portfolioID
		ids := []string{ensureID(&e.ID, "education")}

		for j := range e.Tasks {
			t := &e.Tasks[j]
			if t.Task = strings.TrimSpace(t.Task); t.Task == "" {
				return nil, fmt.Errorf("tasks[%d]: task is required", j)
			}
			t.EducationID = e.ID
			ids = append(ids, ensureID(&t.ID, "education-task"))
		}
		return ids, nil
	}, func(tx *gorm.DB, portfolioID string) error {
		educationIDs := tx.Model(&models.Education{}).Select("id").Where("portfolio_id = ?", portfolioID)
		return tx.Where("education_id IN (?)", educationIDs).Delete(&models.EducationTask{}).Error
	})
}

func HandleAdminReplaceHobbies(c *gin.Context) {
	replacePortfolioList(c, "hobbies", func(h *models.Hobby, portfolioID string) ([]string, error) {
		if h.Hobby = strings.TrimSpace(h.Hobby); h.Hobby == "" {
			return nil, errors.New("hobby is required")
		}
		h.PortfolioID = portfolioID
		return []string{ensureID(&h.ID, "hobby")}, nil
	}, nil)
}

func queryID(c *gin.Context, key string) (uint, bool) {
	id, err := strconv.ParseUint(c.Query(key), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing " + key})
		return 0, false
	}
	return uint(id), true
}

func currentUser(c *gin.Context) *models.User {
	user, _ := c.Get("user")
	u, _ := user.(*models.User)
	return u
}

func newID(prefix string) string {
	b := make([]byte, 8)
	rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}

func loadPortfolio(c *gin.Context) (*models.Portfolio, bool) {
	var p models.Portfolio
	if err := portfolio.Data.GetDB().First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch portfolio"})
		}
		return nil, false
	}
	return &p, true
}

func savePortfolio(c *gin.Context, p *models.Portfolio) {
	if err := portfolio.Data.GetDB().Omit(clause.Associations).Save(p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func replacePortfolioList[T any](
	c *gin.Context,
	name string,
	prepare func(item *T, portfolioID string) ([]string, error),
	beforeDelete func(tx *gorm.DB, portfolioID string) error,
) {
	p, ok := loadPortfolio(c)
	if !ok {
		return
	}

	var items []T
	if err := c.ShouldBindJSON(&items); err != nil || items == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected a JSON array of " + name})
		return
	}

	seen := make(map[string]bool)
	for i := range items {
		ids, err := prepare(&items[i], p.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s[%d]: %v", name, i, err)})
			return
		}
		for _, id := range ids {
			if seen[id] {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s[%d]: duplicate id %q", name, i, id)})
				return
			}
			seen[id] = true
		}
	}

	err := portfolio.Data.GetDB().Transaction(func(tx *gorm.DB) error {
		if beforeDelete != nil {
			if err := beforeDelete(tx, p.ID); err != nil {
				return err
			}
		}
		if err := tx.Where("portfolio_id = ?", p.ID).Delete(new(T)).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Omit("Portfolio").Create(&items).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update " + name})
		return
	}
	c.JSON(http.StatusOK, items)
}

func ensureID(id *string, prefix string) string {
	*id = strings.TrimSpace(*id)
	if *id == "" {
		*id = newID(prefix)
	}
	return *id
}
