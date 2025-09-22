package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForUsers(t *testing.T) *gorm.DB {
	t.Helper()
	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	db.DataBase = d
	if err := d.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate user model: %v", err)
	}
	return d
}

func TestHandleAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	prev := os.Getenv("SECRET_TOKEN")
	defer os.Setenv("SECRET_TOKEN", prev)
	os.Setenv("SECRET_TOKEN", "test-secret")
	d := setupTestDBForUsers(t)

	plain := "password123"
	hashed := internal.Encrypt(plain)
	user := models.User{Name: "bob", Email: "bob@example.com", Password: hashed}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	router := gin.Default()
	router.POST("/v1/login", HandleAuth)

	payload := map[string]string{"name": "bob", "password": plain}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/v1/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["message"] != "Login successful" {
		t.Fatalf("unexpected message: %v", resp["message"])
	}
}
