package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ride-sharing-app/domain/model"
	"github.com/yourusername/ride-sharing-app/domain/repository"
)

// RideService handles ride-related business logic
type RideService struct {
	rideRepo repository.RideRepository
	userRepo repository.UserRepository
}

// NewRideService creates a new RideService
func NewRideService(rideRepo repository.RideRepository, userRepo repository.UserRepository) *RideService {
	return &RideService{
		rideRepo: rideRepo,
		userRepo: userRepo,
	}
}

// CreateRideOffer creates a new ride offer
func (s *RideService) CreateRideOffer(
	driverID uuid.UUID,
	startLat, startLng float64,
	startAddress string,
	endLat, endLng float64,
	endAddress string,
	departureTime time.Time,
	availableSeats int,
	pricePerSeat float64,
	allowedDetourKm float64,
) (*model.RideOffer, error) {
	// Validate driver
	driverProfile, err := s.userRepo.GetDriverProfile(driverID)
	if err != nil {
		return nil, err
	}
	if driverProfile == nil {
		return nil, errors.New("driver profile not found")
	}

	// Check if the available seats is valid
	if availableSeats <= 0 || availableSeats > driverProfile.NumSeats {
		return nil, errors.New("invalid number of available seats")
	}

	// Check if departure time is in the future
	if departureTime.Before(time.Now()) {
		return nil, errors.New("departure time must be in the future")
	}

	// Create new ride offer
	offer := &model.RideOffer{
		DriverID: driverID,
		StartLocation: model.Location{
			Latitude:  startLat,
			Longitude: startLng,
			Address:   startAddress,
		},
		EndLocation: model.Location{
			Latitude:  endLat,
			Longitude: endLng,
			Address:   endAddress,
		},
		DepartureTime:   departureTime,
		AvailableSeats:  availableSeats,
		PricePerSeat:    pricePerSeat,
		AllowedDetourKm: allowedDetourKm,
		Status:          model.StatusPending,
	}

	// Save offer to database
	if err := s.rideRepo.CreateRideOffer(offer); err != nil {
		return nil, err
	}

	// Find potential matches for this offer
	go s.findMatchesForOffer(offer.ID)

	return offer, nil
}

// CreateRideRequest creates a new ride request
func (s *RideService) CreateRideRequest(
	passengerID uuid.UUID,
	startLat, startLng float64,
	startAddress string,
	endLat, endLng float64,
	endAddress string,
	departureTime time.Time,
	numPassengers int,
	maxPrice float64,
) (*model.RideRequest, error) {
	// Validate passenger
	passenger, err := s.userRepo.FindByID(passengerID)
	if err != nil {
		return nil, err
	}
	if passenger == nil {
		return nil, errors.New("passenger not found")
	}

	// Check if the number of passengers is valid
	if numPassengers <= 0 {
		return nil, errors.New("invalid number of passengers")
	}

	// Check if departure time is in the future
	if departureTime.Before(time.Now()) {
		return nil, errors.New("departure time must be in the future")
	}

	// Create new ride request
	request := &model.RideRequest{
		PassengerID: passengerID,
		StartLocation: model.Location{
			Latitude:  startLat,
			Longitude: startLng,
			Address:   startAddress,
		},
		EndLocation: model.Location{
			Latitude:  endLat,
			Longitude: endLng,
			Address:   endAddress,
		},
		DepartureTime: departureTime,
		NumPassengers: numPassengers,
		MaxPrice:      maxPrice,
		Status:        model.StatusPending,
	}

	// Save request to database
	if err := s.rideRepo.CreateRideRequest(request); err != nil {
		return nil, err
	}

	// Find potential matches for this request
	go s.findMatchesForRequest(request.ID)

	return request, nil
}

// GetRideOffersByDriver retrieves all ride offers by a specific driver
func (s *RideService) GetRideOffersByDriver(driverID uuid.UUID) ([]model.RideOffer, error) {
	return s.rideRepo.FindRideOffersByDriverID(driverID)
}

// GetRideRequestsByPassenger retrieves all ride requests by a specific passenger
func (s *RideService) GetRideRequestsByPassenger(passengerID uuid.UUID) ([]model.RideRequest, error) {
	return s.rideRepo.FindRideRequestsByPassengerID(passengerID)
}

// CalculateDistanceBetweenPoints calculates the distance between two geographical points
func (s *RideService) CalculateDistanceBetweenPoints(lat1, lng1, lat2, lng2 float64) float64 {
	// Earth radius in kilometers
	const earthRadius = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlng := lng2Rad - lng1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

// findMatchesForOffer finds and creates potential matches for a ride offer
func (s *RideService) findMatchesForOffer(offerID uuid.UUID) error {
	offer, err := s.rideRepo.FindRideOfferByID(offerID)
	if err != nil {
		return err
	}
	if offer == nil {
		return errors.New("ride offer not found")
	}

	potentialRequests, err := s.rideRepo.FindPotentialMatches(offerID)
	if err != nil {
		return err
	}

	for _, request := range potentialRequests {
		// Calculate match score based on route proximity, time, etc.
		matchScore := s.calculateMatchScore(offer, &request)

		// If match score is good enough, create a match
		if matchScore > 0.6 {
			match := &model.RideMatch{
				RideOfferID:   offer.ID,
				RideRequestID: request.ID,
				Status:        model.StatusMatched,
				MatchScore:    matchScore,
				Price:         offer.PricePerSeat * float64(request.NumPassengers),
			}

			if err := s.rideRepo.CreateRideMatch(match); err != nil {
				// Log error but continue processing other matches
				continue
			}

			// Update offer and request statuses
			offer.Status = model.StatusMatched
			request.Status = model.StatusMatched

			s.rideRepo.UpdateRideOffer(offer)
			s.rideRepo.UpdateRideRequest(&request)
		}
	}

	return nil
}

// findMatchesForRequest finds and creates potential matches for a ride request
func (s *RideService) findMatchesForRequest(requestID uuid.UUID) error {
	request, err := s.rideRepo.FindRideRequestByID(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("ride request not found")
	}

	potentialOffers, err := s.rideRepo.FindPotentialOffers(requestID)
	if err != nil {
		return err
	}

	for _, offer := range potentialOffers {
		// Calculate match score based on route proximity, time, etc.
		matchScore := s.calculateMatchScore(&offer, request)

		// If match score is good enough, create a match
		if matchScore > 0.6 {
			match := &model.RideMatch{
				RideOfferID:   offer.ID,
				RideRequestID: request.ID,
				Status:        model.StatusMatched,
				MatchScore:    matchScore,
				Price:         offer.PricePerSeat * float64(request.NumPassengers),
			}

			if err := s.rideRepo.CreateRideMatch(match); err != nil {
				// Log error but continue processing other matches
				continue
			}

			// Update offer and request statuses
			offer.Status = model.StatusMatched
			request.Status = model.StatusMatched

			s.rideRepo.UpdateRideOffer(&offer)
			s.rideRepo.UpdateRideRequest(request)
		}
	}

	return nil
}

// calculateMatchScore calculates a matching score between an offer and a request
// A higher score means a better match
func (s *RideService) calculateMatchScore(offer *model.RideOffer, request *model.RideRequest) float64 {
	// Check if there are enough seats
	if offer.AvailableSeats < request.NumPassengers {
		return 0
	}

	// Calculate distance from passenger pickup point to driver start point
	pickupDistance := s.CalculateDistanceBetweenPoints(
		request.StartLocation.Latitude,
		request.StartLocation.Longitude,
		offer.StartLocation.Latitude,
		offer.StartLocation.Longitude,
	)

	// Calculate distance from passenger drop-off point to driver end point
	dropoffDistance := s.CalculateDistanceBetweenPoints(
		request.EndLocation.Latitude,
		request.EndLocation.Longitude,
		offer.EndLocation.Latitude,
		offer.EndLocation.Longitude,
	)

	// If the detour is too great, it's not a good match
	if pickupDistance > offer.AllowedDetourKm || dropoffDistance > offer.AllowedDetourKm {
		return 0
	}

	// Calculate time difference in minutes
	timeDiff := math.Abs(float64(offer.DepartureTime.Sub(request.DepartureTime).Minutes()))

	// Calculate price compatibility (0-1)
	pricePerSeat := offer.PricePerSeat
	maxPricePerSeat := request.MaxPrice / float64(request.NumPassengers)
	priceCompat := 1.0
	if pricePerSeat > maxPricePerSeat {
		priceCompat = 0
	}

	// Calculate time compatibility (0-1)
	timeCompat := math.Max(0, 1-timeDiff/30) // 30 minutes is the max time difference

	// Calculate location compatibility (0-1)
	locationCompat := math.Max(0, 1-(pickupDistance+dropoffDistance)/(2*offer.AllowedDetourKm))

	// Weight factors for the final score
	const (
		priceWeight    = 0.3
		timeWeight     = 0.3
		locationWeight = 0.4
	)

	// Calculate weighted average
	return priceCompat*priceWeight + timeCompat*timeWeight + locationCompat*locationWeight
}

// ConfirmMatch confirms a ride match
func (s *RideService) ConfirmMatch(matchID uuid.UUID, userID uuid.UUID, isDriver bool) error {
	match, err := s.rideRepo.FindRideMatchByID(matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found")
	}

	// Verify user is involved in this match
	if isDriver {
		offer, err := s.rideRepo.FindRideOfferByID(match.RideOfferID)
		if err != nil {
			return err
		}
		if offer == nil || offer.DriverID != userID {
			return errors.New("unauthorized: user is not the driver for this match")
		}
	} else {
		request, err := s.rideRepo.FindRideRequestByID(match.RideRequestID)
		if err != nil {
			return err
		}
		if request == nil || request.PassengerID != userID {
			return errors.New("unauthorized: user is not the passenger for this match")
		}
	}

	// Update match status
	if match.Status != model.StatusMatched {
		return errors.New("match cannot be confirmed in its current state")
	}

	match.Status = model.StatusConfirmed

	// Update offer and request statuses
	offer, err := s.rideRepo.FindRideOfferByID(match.RideOfferID)
	if err != nil {
		return err
	}
	request, err := s.rideRepo.FindRideRequestByID(match.RideRequestID)
	if err != nil {
		return err
	}

	offer.Status = model.StatusConfirmed
	request.Status = model.StatusConfirmed

	// Update available seats
	offer.AvailableSeats -= request.NumPassengers

	if err := s.rideRepo.UpdateRideMatch(match); err != nil {
		return err
	}
	if err := s.rideRepo.UpdateRideOffer(offer); err != nil {
		return err
	}
	if err := s.rideRepo.UpdateRideRequest(request); err != nil {
		return err
	}

	return nil
}
