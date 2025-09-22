package views

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"github.com/gin-gonic/gin"
)

func TestHandleGetBlogComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d := SetupInMemoryDB(t)

	user := models.User{Name: "cuser", Email: "c@e.com", Password: "x"}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	post := models.BlogPost{Title: "pc", Content: "cc", UserID: user.ID, IsPublished: true}
	if err := d.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	if err := d.Create(&models.BlogComment{BlogPostID: post.ID, UserID: user.ID, Comment: "first"}).Error; err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	router := gin.Default()
	router.GET("/handle/blog/comments", func(c *gin.Context) {
		HandleGetBlogComments(c)
	})

	req := httptest.NewRequest("GET", "/handle/blog/comments?post_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}
