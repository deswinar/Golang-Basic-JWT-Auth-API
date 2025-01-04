package models

import (
	"gorm.io/gorm"
)

// Arena represents a sport arena where users can book fields
type Arena struct {
	gorm.Model
	Name        string    `gorm:"not null" json:"name"`
	Location    string    `gorm:"not null" json:"location"`
	OwnerID     uint      `gorm:"not null" json:"owner_id"` // Foreign key for User (Owner)
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner"`
	Fields      []Field   `gorm:"foreignKey:ArenaID" json:"fields"`
	Bookings    []Booking `gorm:"foreignKey:ArenaID" json:"bookings"`
}

// Field represents a sport field inside an arena
type Field struct {
	gorm.Model
	ArenaID    uint   `gorm:"not null" json:"arena_id"` // Foreign key for Arena
	FieldName  string `gorm:"not null" json:"field_name"`
	SportType  string `gorm:"not null" json:"sport_type"` // e.g., Soccer, Basketball, etc.
	PricePerHr float64 `json:"price_per_hr"` // Price per hour for the field
}