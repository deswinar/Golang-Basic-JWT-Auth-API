package middleware

import (
	"net/http"
	"sparring-backend/auth"
	"sparring-backend/internal/database"
	"sparring-backend/internal/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and handles rotating secrets
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Validate token using rotating secrets
		userID, claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Optionally, check claims (e.g., user_id, roles) for more validation
		if _, ok := claims["user_id"].(float64); !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Retrieve the user from the database using the userID
		var user models.User
		err = database.DB.First(&user, userID).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Save the user data to the context
		c.Set("user", user) // Store full user model (not just userID)
		c.Set("claims", claims)
		c.Next()
	}
}

type RateLimiter struct {
	// Store the timestamps of requests for each IP address.
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int           // Max requests allowed in a time window
	window   time.Duration // Time window for rate limit
}

// NewRateLimiter creates a new RateLimiter with a specified limit and time window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Limit is a middleware that limits the number of requests to an endpoint.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the IP address of the requester
		ip := c.ClientIP()

		// Lock to ensure thread-safety when updating the requests map
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Get current time
		now := time.Now()

		// Initialize the IP entry if it doesn't exist
		if _, exists := rl.requests[ip]; !exists {
			rl.requests[ip] = []time.Time{}
		}

		// Remove old requests that are outside the time window
		validRequests := []time.Time{}
		for _, t := range rl.requests[ip] {
			if now.Sub(t) < rl.window {
				validRequests = append(validRequests, t)
			}
		}
		rl.requests[ip] = validRequests

		// Check if the number of requests exceeds the limit
		if len(rl.requests[ip]) >= rl.limit {
			// Too many requests, send a 429 response
			c.JSON(429, gin.H{
				"message": "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}

		// Add the current request timestamp to the list
		rl.requests[ip] = append(rl.requests[ip], now)

		// Continue to the next handler
		c.Next()
	}
}
