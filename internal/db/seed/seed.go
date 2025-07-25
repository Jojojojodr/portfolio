package seed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Jojojojodr/portfolio/internal"
	"github.com/Jojojojodr/portfolio/internal/db/models"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
    seedDir := "database/seeds"

    if _, err := os.Stat(seedDir); os.IsNotExist(err) {
        log.Printf("Seeds directory %s does not exist, skipping seeding", seedDir)
        return nil
    }

    log.Println("Starting database seeding...")

    if err := seedUsers(db, filepath.Join(seedDir, "users.json")); err != nil {
        return fmt.Errorf("failed to seed users: %v", err)
    }

    if err := seedFromJSON(db, filepath.Join(seedDir, "blog_posts.json"), &[]models.BlogPost{}); err != nil {
        return fmt.Errorf("failed to seed blog posts: %v", err)
    }

    if err := seedFromJSON(db, filepath.Join(seedDir, "blog_comments.json"), &[]models.BlogComment{}); err != nil {
        return fmt.Errorf("failed to seed blog comments: %v", err)
    }

    if err := seedFromJSON(db, filepath.Join(seedDir, "post_likes.json"), &[]models.PostLike{}); err != nil {
        return fmt.Errorf("failed to seed post likes: %v", err)
    }

    if err := seedFromJSON(db, filepath.Join(seedDir, "comment_likes.json"), &[]models.CommentLike{}); err != nil {
        return fmt.Errorf("failed to seed comment likes: %v", err)
    }

    log.Println("Database seeding completed successfully")
    return nil
}

func seedUsers(db *gorm.DB, filePath string) error {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        log.Printf("Seed file %s does not exist, skipping", filePath)
        return nil
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %v", filePath, err)
    }

    var users []models.User
    if err := json.Unmarshal(data, &users); err != nil {
        return fmt.Errorf("failed to unmarshal JSON from %s: %v", filePath, err)
    }

    for i := range users {
        if users[i].Password != "" {
            hashedPassword := internal.Encrypt(users[i].Password)
            users[i].Password = hashedPassword
        }
    }

    if len(users) > 0 {
        if err := db.Create(&users).Error; err != nil {
            return fmt.Errorf("failed to insert users: %v", err)
        }
    }

    log.Printf("Successfully seeded %d users from %s", len(users), filePath)
    return nil
}

func seedFromJSON(db *gorm.DB, filePath string, model interface{}) error {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        log.Printf("Seed file %s does not exist, skipping", filePath)
        return nil
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %v", filePath, err)
    }

    if err := json.Unmarshal(data, model); err != nil {
        return fmt.Errorf("failed to unmarshal JSON from %s: %v", filePath, err)
    }

    if err := db.Create(model).Error; err != nil {
        return fmt.Errorf("failed to insert data from %s: %v", filePath, err)
    }

    log.Printf("Successfully seeded data from %s", filePath)
    return nil
}