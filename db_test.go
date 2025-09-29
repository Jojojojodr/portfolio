package portfolio

import (
	"os"
	"testing"

	"github.com/Jojojojodr/portfolio/internal/db/models"
)

func TestConnectSQLite_InMemory(t *testing.T) {
	prevType := os.Getenv("DB_TYPE")
	prevPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_TYPE", prevType)
	defer os.Setenv("DB_PATH", prevPath)

	os.Setenv("DB_TYPE", "sqlite")
	os.Setenv("DB_PATH", "file::memory:?cache=shared")

	gdb := ConnectDB(prevType)
	if gdb == nil {
		t.Fatalf("ConnectDB returned nil")
	}

	if ok := gdb.Migrator().HasTable(&models.User{}); !ok {
		t.Fatalf("expected users table to exist after migration")
	}
	if ok := gdb.Migrator().HasTable(&models.BlogPost{}); !ok {
		t.Fatalf("expected blog_posts table to exist after migration")
	}
	if ok := gdb.Migrator().HasTable(&models.BlogComment{}); !ok {
		t.Fatalf("expected blog_comments table to exist after migration")
	}
}

func Test_connectSQLite_DirectInMemory(t *testing.T) {
	prevPath := os.Getenv("DB_PATH")
	defer os.Setenv("DB_PATH", prevPath)
	os.Setenv("DB_PATH", "file::memory:?cache=shared")

	d, err := connectSQLite()
	if err != nil {
		t.Fatalf("connectSQLite returned error: %v", err)
	}
	if d == nil {
		t.Fatalf("connectSQLite returned nil DB")
	}
}
