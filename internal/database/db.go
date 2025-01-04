package database

import (
	"log"
	"os"
	"sparring-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

// InitDB initializes the database connection and automates migrations for models
func InitDB() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
		return err
	}

	// Retrieve the database connection details from environment variables
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_CHARSET := os.Getenv("DB_CHARSET")
	DB_LOC := os.Getenv("DB_LOC")

	// Build the database connection string
	dsn := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=" + DB_CHARSET + "&parseTime=True&loc=" + DB_LOC

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return err
	}

	// Automigrate the models
	err = DB.AutoMigrate(&models.User{}, &models.Arena{}, &models.Booking{}) // Add more models here
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}

	log.Println("Database connected and migrated successfully.")
	return nil
}
