package views

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"github.com/gin-gonic/gin"
)

func TestLogin_ViewHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d := SetupInMemoryDB(t)

	plain := "mypw"
	hashed := internal.Encrypt(plain)
	user := models.User{Name: "vuser", Email: "v@e.com", Password: hashed}
	if err := d.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	
	prev := os.Getenv("SECRET_TOKEN")
	defer os.Setenv("SECRET_TOKEN", prev)
	os.Setenv("SECRET_TOKEN", "test-secret-views")

	form := url.Values{}
	form.Add("name", "vuser")
	form.Add("password", plain)

	router := gin.Default()
	router.POST("/handle/auth/login", func(c *gin.Context) { Login(c) })

	req := httptest.NewRequest("POST", "/handle/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_PostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = SetupInMemoryDB(t)

	form := url.Values{}
	form.Add("name", "newuser")
	form.Add("email", "new@e.com")
	form.Add("password", "newpass")

	router := gin.Default()
	router.POST("/handle/register", func(c *gin.Context) { HandleRegisterPost(c) })

	req := httptest.NewRequest("POST", "/handle/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303 got %d: %s", w.Code, w.Body.String())
	}
}
