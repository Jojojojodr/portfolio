package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleHealth(c *gin.Context) {
	// Health Check
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}