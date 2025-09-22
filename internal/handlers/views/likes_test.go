package views

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"github.com/gin-gonic/gin"
)

func TestTogglePostLike(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d := SetupInMemoryDB(t)

	user := models.User{Name: "liker", Email: "l@e.com", Password: "x"}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	post := models.BlogPost{Title: "p", Content: "c", UserID: user.ID, IsPublished: true}
	if err := d.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	router := gin.Default()
	router.POST("/handle/like/post/:postId", func(c *gin.Context) {
		c.Set("user", &user)
		TogglePostLike(c)
	})

	req := httptest.NewRequest("POST", "/handle/like/post/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestToggleCommentLike(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d := SetupInMemoryDB(t)

	user := models.User{Name: "liker2", Email: "l2@e.com", Password: "x"}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	post := models.BlogPost{Title: "p2", Content: "c2", UserID: user.ID, IsPublished: true}
	if err := d.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	comment := models.BlogComment{BlogPostID: post.ID, UserID: user.ID, Comment: "hey"}
	if err := d.Create(&comment).Error; err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	router := gin.Default()
	router.POST("/handle/like/comment/:commentId", func(c *gin.Context) {
		c.Set("user", &user)
		ToggleCommentLike(c)
	})

	req := httptest.NewRequest("POST", "/handle/like/comment/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}
