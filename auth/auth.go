package auth

import (
	"crypto/rand" // for cryptographic random generation
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Load secrets from environment variables
var jwtSecrets = struct {
	ActiveSecret   string
	ExpiredSecrets []string
}{
	ActiveSecret:   "",
	ExpiredSecrets: []string{},
}

func init() {
	// Load environment variables
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load active JWT secret from environment variables
	jwtSecrets.ActiveSecret = os.Getenv("JWT_ACTIVE_SECRET")

	// Load expired secrets from environment variables (if any)
	jwtSecrets.ExpiredSecrets = append(jwtSecrets.ExpiredSecrets, os.Getenv("JWT_OLD_SECRET_1"))
	jwtSecrets.ExpiredSecrets = append(jwtSecrets.ExpiredSecrets, os.Getenv("JWT_OLD_SECRET_2"))
}

// Global variable or function to fetch the active secret key
func GetActiveSecret() string {
	// Logic to return the current active secret, considering rotation
	// For example, fetching it from a map of secrets or reading from environment variables
	return jwtSecrets.ActiveSecret // or rotate this logic
}

// GenerateToken generates a JWT token for the given user ID.
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := GetActiveSecret()               // Get active secret based on rotation logic
	return token.SignedString([]byte(secret)) // Sign using the active secret
}

// GenerateRefreshToken generates a refresh token for the given user ID.
func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days expiry for refresh token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := GetActiveSecret()               // Get active secret based on rotation logic
	return token.SignedString([]byte(secret)) // Sign using the active secret
}

// ValidateToken validates the JWT token and returns the user ID if valid.
func ValidateToken(tokenString string) (uint, jwt.MapClaims, error) {
	// Strip the "Bearer " prefix if present
	tokenString = stripBearerPrefix(tokenString)
	if tokenString == "" {
		return 0, nil, errors.New("token is missing or invalid format")
	}

	// Parse the token using the active secret first
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecrets.ActiveSecret), nil
	})
	if err == nil && token.Valid {
		// Token is valid with the active secret
		return extractUserIDFromClaims(token.Claims)
	}

	// If token validation fails with the active secret, try expired secrets
	for _, secret := range jwtSecrets.ExpiredSecrets {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err == nil && token.Valid {
			// Token is valid with an expired secret
			return extractUserIDFromClaims(token.Claims)
		}
	}

	return 0, nil, errors.New("invalid token or expired secret")
}

// extractUserIDFromClaims extracts the user ID from the token claims and returns the claims.
func extractUserIDFromClaims(claims jwt.Claims) (uint, jwt.MapClaims, error) {
	// Validate the claims and extract the user_id
	if mapClaims, ok := claims.(jwt.MapClaims); ok {
		userID, ok := mapClaims["user_id"].(float64)
		if !ok {
			return 0, nil, errors.New("user_id not found in token")
		}
		return uint(userID), mapClaims, nil
	}
	return 0, nil, errors.New("invalid claims")
}

// stripBearerPrefix removes the "Bearer " prefix from the token string.
func stripBearerPrefix(tokenString string) string {
	if strings.HasPrefix(tokenString, "Bearer ") {
		return tokenString[7:]
	}
	return tokenString
}

// RotateSecret rotates the JWT secret by saving the current active secret as expired and generating a new one.
func RotateSecret() {
	// Save the current active secret to expired secrets
	jwtSecrets.ExpiredSecrets = append(jwtSecrets.ExpiredSecrets, jwtSecrets.ActiveSecret)

	// Limit the number of expired secrets stored (e.g., keep only 3 recent secrets)
	if len(jwtSecrets.ExpiredSecrets) > 3 {
		jwtSecrets.ExpiredSecrets = jwtSecrets.ExpiredSecrets[1:]
	}

	// Generate a new secret and store it in the environment (environment variables can be updated manually or via CI/CD)
	newSecret := generateRandomSecret(32)
	jwtSecrets.ActiveSecret = newSecret

	// Example: Update the environment variable for the new active secret
	// os.Setenv("JWT_ACTIVE_SECRET", newSecret) // Update in a secure way
	log.Println("JWT Secret rotated successfully. New secret generated.")
}

// generateRandomSecret generates a cryptographically secure random JWT secret.
func generateRandomSecret(length int) string {
	secret := make([]byte, length)
	_, err := rand.Read(secret) // Use crypto/rand for secure randomness
	if err != nil {
		log.Fatal("Error generating random secret: ", err)
	}

	// Convert to string using a base64 encoding
	return string(secret)
}

func StartSecretRotation(interval time.Duration) {
    go func() {
        for {
            time.Sleep(interval)
            RotateSecret()
            log.Println("JWT secret rotated")
        }
    }()
}
