package portfolio

import (
	"fmt"
	"net"
	"strings"

	"github.com/Jojojojodr/portfolio/config"
	"github.com/Jojojojodr/portfolio/database/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SecretToken string

func Encrypt(key string) string {
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error encrypting key: %v\n", err)
		return ""
	}
	return string(hashedKey)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsDatabaseEmpty(database *gorm.DB) bool {
    var userCount int64 
    database.Model(&models.User{}).Count(&userCount)
    return userCount == 0
}

func SetSecretToken(token string) {
	SecretToken = token
}

func PublicURL(c *gin.Context, path string) string {
	if !isAdminRequest(c) {
		return path
	}
	return siteURL(c, config.AppConfig.Server.Port, path)
}

func AdminURL(c *gin.Context, path string) string {
	if isAdminRequest(c) {
		return path
	}
	return siteURL(c, config.AppConfig.Server.AdminPort, path)
}

func isAdminRequest(c *gin.Context) bool {
	port := requestPort(c)
	return port != "" && port == config.AppConfig.Server.AdminPort
}

func requestPort(c *gin.Context) string {
	_, port, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		return ""
	}
	return port
}

func siteURL(c *gin.Context, port, path string) string {
	if port == "" {
		return path
	}
	host := c.Request.Host
	if _, currentPort, err := net.SplitHostPort(host); err == nil {
		host = strings.TrimSuffix(host, ":"+currentPort)
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + host + ":" + port + path
}
