package views

import (
	"testing"

	"github.com/Jojojojodr/portfolio/internal/db"
	"github.com/Jojojojodr/portfolio/internal/db/models"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupInMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	d, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	db.DataBase = d

	if err := d.AutoMigrate(&models.User{}, &models.BlogPost{}, &models.BlogComment{}, &models.PostLike{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	return d
}
