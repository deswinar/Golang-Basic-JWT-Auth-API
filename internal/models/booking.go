package models

import (
	"time"

	"gorm.io/gorm"
)

// Booking represents a user's booking for a field in an arena
type Booking struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"` // Foreign key for User
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	ArenaID     uint      `gorm:"not null" json:"arena_id"` // Foreign key for Arena
	Arena       Arena     `gorm:"foreignKey:ArenaID" json:"arena"`
	FieldID     uint      `gorm:"not null" json:"field_id"` // Foreign key for Field
	Field       Field     `gorm:"foreignKey:FieldID" json:"field"`
	BookingTime time.Time `gorm:"not null" json:"booking_time"` // Booking time for the field
	Duration    int       `gorm:"not null" json:"duration"`     // Duration in hours
	TotalAmount float64   `json:"total_amount"`                 // Total amount to be paid
	Status      string    `gorm:"not null" json:"status"`       // "pending", "confirmed", "cancelled", etc.
}
