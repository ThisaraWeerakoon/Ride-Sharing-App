package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the role of a user in the system
type UserRole string

const (
	// RolePassenger represents a user who can request rides
	RolePassenger UserRole = "passenger"
	// RoleDriver represents a user who can offer rides
	RoleDriver UserRole = "driver"
	// RoleBoth represents a user who can both request and offer rides
	RoleBoth UserRole = "both"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate generates a UUID for new users before creating them
func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// DriverProfile contains additional information for users with driver role
type DriverProfile struct {
	UserID        uuid.UUID `json:"-" gorm:"primaryKey;type:uuid"`
	User          User      `json:"-" gorm:"foreignKey:UserID"`
	LicenseNo     string    `json:"license_no" gorm:"not null"`
	CarModel      string    `json:"car_model" gorm:"not null"`
	CarPlateNo    string    `json:"car_plate_no" gorm:"not null"`
	NumSeats      int       `json:"num_seats" gorm:"not null"`
	AverageRating float32   `json:"avg_rating" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
