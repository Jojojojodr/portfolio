package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForBlogCreate(t *testing.T) *gorm.DB {
	t.Helper()
	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	db.DataBase = d
	if err := d.AutoMigrate(&models.User{}, &models.BlogPost{}); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}
	return d
}

func TestCreateBlogPost_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d := setupTestDBForBlogCreate(t)

	user := models.User{Name: "alice", Email: "alice@example.com", Password: "x"}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	router := gin.Default()
	router.POST("/v1/blog/create", CreateBlogPost)

	payload := map[string]interface{}{
		"title":        "New Post",
		"content":      "This is a test",
		"user_id":      user.ID,
		"is_published": true,
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/v1/blog/create", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created models.BlogPost
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if created.Title != "New Post" || created.UserID != user.ID {
		t.Fatalf("created post mismatch: %+v", created)
	}
}
