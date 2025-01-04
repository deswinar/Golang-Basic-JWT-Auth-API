package main

import (
	"log"
	"net/http"
	"sparring-backend/auth"
	"sparring-backend/handlers"
	"sparring-backend/internal/database"
	"sparring-backend/internal/models"
	"sparring-backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
		return
	}

	// Rotate secret every 30 days
	auth.StartSecretRotation(30 * 24 * time.Hour)

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default()) // Use the default CORS policy

	// Initialize the rate limiter (limit to 5 requests per minute)
	rateLimiter := middleware.NewRateLimiter(5, 1*time.Minute)

	// Define routes for user, arena, and booking
	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", rateLimiter.Limit(), handlers.LoginUser)                  // Apply rate limiting to login
	router.POST("/refresh-token", rateLimiter.Limit(), handlers.RefreshAccessToken) // Apply rate limiting to refresh-token
	router.POST("/logout", handlers.LogoutUser)
	router.POST("/arena", handlers.CreateArena)
	router.POST("/booking", handlers.CreateBooking)

	// Add the test endpoint here
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is working!",
		})
	})

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Use authentication middleware

	// Route to fetch current user details
	protected.GET("/protected", func(c *gin.Context) {
		// Retrieve the user from the context
		user, exists := c.Get("user")

		// Check if the user exists in the context and handle the case where it doesn't
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User context not found"})
			return
		}

		// Cast the user to the correct type (User struct)
		currentUser, ok := user.(models.User)
		if !ok {
			// If the cast fails, return an error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}

		// Return user details (name and email)
		c.JSON(200, gin.H{
			"message": "You have access to this protected endpoint",
			"user": gin.H{
				"name":  currentUser.Name,
				"email": currentUser.Email,
			},
		})
	})

	// Start the server
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
