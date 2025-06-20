package model

import (
	"time"

	"github.com/google/uuid"
)

// RideStatus defines the current status of a ride
type RideStatus string

const (
	// StatusPending indicates a ride is waiting for matches
	StatusPending RideStatus = "pending"
	// StatusMatched indicates a ride has been matched
	StatusMatched RideStatus = "matched"
	// StatusConfirmed indicates a ride has been confirmed by all parties
	StatusConfirmed RideStatus = "confirmed"
	// StatusInProgress indicates a ride is currently in progress
	StatusInProgress RideStatus = "in_progress"
	// StatusCompleted indicates a ride has been completed
	StatusCompleted RideStatus = "completed"
	// StatusCancelled indicates a ride has been cancelled
	StatusCancelled RideStatus = "cancelled"
)

// Location represents a geographical point
type Location struct {
	Latitude  float64 `json:"lat" gorm:"not null"`
	Longitude float64 `json:"lng" gorm:"not null"`
	Address   string  `json:"address" gorm:"not null"`
}

// RideOffer represents a ride offered by a driver
type RideOffer struct {
	ID              uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	DriverID        uuid.UUID  `json:"driver_id" gorm:"type:uuid;not null"`
	Driver          User       `json:"-" gorm:"foreignKey:DriverID"`
	StartLocation   Location   `json:"start_location" gorm:"embedded;embeddedPrefix:start_"`
	EndLocation     Location   `json:"end_location" gorm:"embedded;embeddedPrefix:end_"`
	DepartureTime   time.Time  `json:"departure_time" gorm:"not null"`
	AvailableSeats  int        `json:"available_seats" gorm:"not null"`
	Status          RideStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	PricePerSeat    float64    `json:"price_per_seat" gorm:"not null"`
	AllowedDetourKm float64    `json:"allowed_detour_km" gorm:"default:5"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate generates a UUID for new ride offers before creating them
func (r *RideOffer) BeforeCreate() error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RideRequest represents a ride requested by a passenger
type RideRequest struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	PassengerID   uuid.UUID  `json:"passenger_id" gorm:"type:uuid;not null"`
	Passenger     User       `json:"-" gorm:"foreignKey:PassengerID"`
	StartLocation Location   `json:"start_location" gorm:"embedded;embeddedPrefix:start_"`
	EndLocation   Location   `json:"end_location" gorm:"embedded;embeddedPrefix:end_"`
	DepartureTime time.Time  `json:"departure_time" gorm:"not null"`
	NumPassengers int        `json:"num_passengers" gorm:"not null;default:1"`
	Status        RideStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	MaxPrice      float64    `json:"max_price" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate generates a UUID for new ride requests before creating them
func (r *RideRequest) BeforeCreate() error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RideMatch represents a match between a ride offer and request
type RideMatch struct {
	ID            uuid.UUID   `json:"id" gorm:"primaryKey;type:uuid"`
	RideOfferID   uuid.UUID   `json:"ride_offer_id" gorm:"type:uuid;not null"`
	RideOffer     RideOffer   `json:"-" gorm:"foreignKey:RideOfferID"`
	RideRequestID uuid.UUID   `json:"ride_request_id" gorm:"type:uuid;not null"`
	RideRequest   RideRequest `json:"-" gorm:"foreignKey:RideRequestID"`
	Status        RideStatus  `json:"status" gorm:"type:varchar(20);default:'matched'"`
	MatchScore    float64     `json:"match_score" gorm:"not null"`
	Price         float64     `json:"price" gorm:"not null"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate generates a UUID for new ride matches before creating them
func (r *RideMatch) BeforeCreate() error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
