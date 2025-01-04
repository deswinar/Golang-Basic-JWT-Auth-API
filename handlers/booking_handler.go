package handlers

import (
	"net/http"
	"sparring-backend/internal/database"
	"sparring-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateBooking handles booking a sports arena
func CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save booking to database
	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
}
