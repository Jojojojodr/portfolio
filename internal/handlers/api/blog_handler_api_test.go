package api

import (
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

func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open in-memory db: %v", err)
    }

    db.DataBase = d

    if err := d.AutoMigrate(&models.BlogPost{}, &models.BlogComment{}, &models.User{}); err != nil {
        t.Fatalf("failed to migrate models: %v", err)
    }

    return d
}

func TestGetBlogPostByID_Handler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    d := setupTestDB(t)

    user := models.User{ID: 1, Name: "testuser", Email: "test@example.com", Password: "secret"}
    if err := d.Create(&user).Error; err != nil {
        t.Fatalf("failed to create user: %v", err)
    }

    post := models.BlogPost{Title: "Hello", Content: "World", UserID: user.ID, IsPublished: true}
    if err := d.Create(&post).Error; err != nil {
        t.Fatalf("failed to create post: %v", err)
    }

    comment := models.BlogComment{Comment: "Nice post", UserID: user.ID, BlogPostID: post.ID}
    if err := d.Create(&comment).Error; err != nil {
        t.Fatalf("failed to create comment: %v", err)
    }

    router := gin.Default()
    router.GET("/v1/blog/post", GetBlogPostByID)

    req := httptest.NewRequest("GET", "/v1/blog/post?id=1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
    }

    var resp struct {
        Post     models.BlogPost      `json:"post"`
        Comments []models.BlogComment `json:"comments"`
    }

    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }

    if resp.Post.ID == 0 || resp.Post.Title != "Hello" {
        t.Fatalf("unexpected post in response: %+v", resp.Post)
    }

    if len(resp.Comments) != 1 || resp.Comments[0].Comment != "Nice post" {
        t.Fatalf("unexpected comments: %+v", resp.Comments)
    }
}