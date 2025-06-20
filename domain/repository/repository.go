package repository

import (
	"github.com/google/uuid"
	"github.com/yourusername/ride-sharing-app/domain/model"
)

// UserRepository defines the contract for user data operations
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	CreateDriverProfile(profile *model.DriverProfile) error
	GetDriverProfile(userID uuid.UUID) (*model.DriverProfile, error)
	UpdateDriverProfile(profile *model.DriverProfile) error
}

// RideRepository defines the contract for ride operations
type RideRepository interface {
	// Ride Offer operations
	CreateRideOffer(offer *model.RideOffer) error
	FindRideOfferByID(id uuid.UUID) (*model.RideOffer, error)
	FindRideOffersByDriverID(driverID uuid.UUID) ([]model.RideOffer, error)
	UpdateRideOffer(offer *model.RideOffer) error
	DeleteRideOffer(id uuid.UUID) error

	// Ride Request operations
	CreateRideRequest(request *model.RideRequest) error
	FindRideRequestByID(id uuid.UUID) (*model.RideRequest, error)
	FindRideRequestsByPassengerID(passengerID uuid.UUID) ([]model.RideRequest, error)
	UpdateRideRequest(request *model.RideRequest) error
	DeleteRideRequest(id uuid.UUID) error

	// Ride Match operations
	CreateRideMatch(match *model.RideMatch) error
	FindRideMatchByID(id uuid.UUID) (*model.RideMatch, error)
	FindRideMatchesByOfferID(offerID uuid.UUID) ([]model.RideMatch, error)
	FindRideMatchesByRequestID(requestID uuid.UUID) ([]model.RideMatch, error)
	UpdateRideMatch(match *model.RideMatch) error
	DeleteRideMatch(id uuid.UUID) error

	// Match finding operations
	FindPotentialMatches(offerID uuid.UUID) ([]model.RideRequest, error)
	FindPotentialOffers(requestID uuid.UUID) ([]model.RideOffer, error)
}
