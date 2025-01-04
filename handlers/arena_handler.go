package handlers

import (
	"net/http"
	"sparring-backend/internal/database"
	"sparring-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateArena handles creating a new arena
func CreateArena(c *gin.Context) {
	var arena models.Arena
	if err := c.ShouldBindJSON(&arena); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save arena to database
	if err := database.DB.Create(&arena).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create arena"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Arena created successfully"})
}
