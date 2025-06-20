package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/yourusername/ride-sharing-app/domain/model"
	repo "github.com/yourusername/ride-sharing-app/domain/repository"
)

// GormRideRepository is an implementation of RideRepository using Gorm
type GormRideRepository struct {
	db *gorm.DB
}

// NewGormRideRepository creates a new GormRideRepository
func NewGormRideRepository(db *gorm.DB) repo.RideRepository {
	return &GormRideRepository{db: db}
}

// CreateRideOffer adds a new ride offer to the database
func (r *GormRideRepository) CreateRideOffer(offer *model.RideOffer) error {
	return r.db.Create(offer).Error
}

// FindRideOfferByID retrieves a ride offer by ID
func (r *GormRideRepository) FindRideOfferByID(id uuid.UUID) (*model.RideOffer, error) {
	var offer model.RideOffer
	if err := r.db.Where("id = ?", id).First(&offer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &offer, nil
}

// FindRideOffersByDriverID retrieves all ride offers by a specific driver
func (r *GormRideRepository) FindRideOffersByDriverID(driverID uuid.UUID) ([]model.RideOffer, error) {
	var offers []model.RideOffer
	if err := r.db.Where("driver_id = ?", driverID).Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

// UpdateRideOffer updates a ride offer in the database
func (r *GormRideRepository) UpdateRideOffer(offer *model.RideOffer) error {
	return r.db.Save(offer).Error
}

// DeleteRideOffer removes a ride offer from the database
func (r *GormRideRepository) DeleteRideOffer(id uuid.UUID) error {
	return r.db.Delete(&model.RideOffer{}, "id = ?", id).Error
}

// CreateRideRequest adds a new ride request to the database
func (r *GormRideRepository) CreateRideRequest(request *model.RideRequest) error {
	return r.db.Create(request).Error
}

// FindRideRequestByID retrieves a ride request by ID
func (r *GormRideRepository) FindRideRequestByID(id uuid.UUID) (*model.RideRequest, error) {
	var request model.RideRequest
	if err := r.db.Where("id = ?", id).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// FindRideRequestsByPassengerID retrieves all ride requests by a specific passenger
func (r *GormRideRepository) FindRideRequestsByPassengerID(passengerID uuid.UUID) ([]model.RideRequest, error) {
	var requests []model.RideRequest
	if err := r.db.Where("passenger_id = ?", passengerID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// UpdateRideRequest updates a ride request in the database
func (r *GormRideRepository) UpdateRideRequest(request *model.RideRequest) error {
	return r.db.Save(request).Error
}

// DeleteRideRequest removes a ride request from the database
func (r *GormRideRepository) DeleteRideRequest(id uuid.UUID) error {
	return r.db.Delete(&model.RideRequest{}, "id = ?", id).Error
}

// CreateRideMatch adds a new ride match to the database
func (r *GormRideRepository) CreateRideMatch(match *model.RideMatch) error {
	return r.db.Create(match).Error
}

// FindRideMatchByID retrieves a ride match by ID
func (r *GormRideRepository) FindRideMatchByID(id uuid.UUID) (*model.RideMatch, error) {
	var match model.RideMatch
	if err := r.db.Where("id = ?", id).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

// FindRideMatchesByOfferID retrieves all matches for a specific ride offer
func (r *GormRideRepository) FindRideMatchesByOfferID(offerID uuid.UUID) ([]model.RideMatch, error) {
	var matches []model.RideMatch
	if err := r.db.Where("ride_offer_id = ?", offerID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// FindRideMatchesByRequestID retrieves all matches for a specific ride request
func (r *GormRideRepository) FindRideMatchesByRequestID(requestID uuid.UUID) ([]model.RideMatch, error) {
	var matches []model.RideMatch
	if err := r.db.Where("ride_request_id = ?", requestID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

// UpdateRideMatch updates a ride match in the database
func (r *GormRideRepository) UpdateRideMatch(match *model.RideMatch) error {
	return r.db.Save(match).Error
}

// DeleteRideMatch removes a ride match from the database
func (r *GormRideRepository) DeleteRideMatch(id uuid.UUID) error {
	return r.db.Delete(&model.RideMatch{}, "id = ?", id).Error
}

// FindPotentialMatches finds potential ride requests that match a ride offer
func (r *GormRideRepository) FindPotentialMatches(offerID uuid.UUID) ([]model.RideRequest, error) {
	var offer model.RideOffer
	if err := r.db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		return nil, err
	}

	// Define the time window for potential matches (e.g., +/- 30 minutes)
	timeWindow := 30 * time.Minute
	startTime := offer.DepartureTime.Add(-timeWindow)
	endTime := offer.DepartureTime.Add(timeWindow)

	var requests []model.RideRequest
	// Find requests within the same timeframe and with compatible locations
	// This is a simplified version, in a real app you would use more sophisticated geospatial queries
	if err := r.db.Where("status = ? AND departure_time BETWEEN ? AND ? AND num_passengers <= ?",
		model.StatusPending, startTime, endTime, offer.AvailableSeats).
		Find(&requests).Error; err != nil {
		return nil, err
	}

	// Further filtering would be done in the service layer
	return requests, nil
}

// FindPotentialOffers finds potential ride offers that match a ride request
func (r *GormRideRepository) FindPotentialOffers(requestID uuid.UUID) ([]model.RideOffer, error) {
	var request model.RideRequest
	if err := r.db.Where("id = ?", requestID).First(&request).Error; err != nil {
		return nil, err
	}

	// Define the time window for potential matches (e.g., +/- 30 minutes)
	timeWindow := 30 * time.Minute
	startTime := request.DepartureTime.Add(-timeWindow)
	endTime := request.DepartureTime.Add(timeWindow)

	var offers []model.RideOffer
	// Find offers within the same timeframe and with sufficient seats
	// This is a simplified version, in a real app you would use more sophisticated geospatial queries
	if err := r.db.Where("status = ? AND departure_time BETWEEN ? AND ? AND available_seats >= ? AND price_per_seat <= ?",
		model.StatusPending, startTime, endTime, request.NumPassengers, request.MaxPrice/float64(request.NumPassengers)).
		Find(&offers).Error; err != nil {
		return nil, err
	}

	// Further filtering would be done in the service layer
	return offers, nil
}
