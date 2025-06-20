package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ride-sharing-app/service"
)

// RideHandler handles ride-related API requests
type RideHandler struct {
	rideService *service.RideService
}

// NewRideHandler creates a new RideHandler
func NewRideHandler(rideService *service.RideService) *RideHandler {
	return &RideHandler{
		rideService: rideService,
	}
}

// CreateRideOfferRequest represents the request format for creating a ride offer
type CreateRideOfferRequest struct {
	StartLocation struct {
		Latitude  float64 `json:"lat" binding:"required"`
		Longitude float64 `json:"lng" binding:"required"`
		Address   string  `json:"address" binding:"required"`
	} `json:"start_location" binding:"required"`

	EndLocation struct {
		Latitude  float64 `json:"lat" binding:"required"`
		Longitude float64 `json:"lng" binding:"required"`
		Address   string  `json:"address" binding:"required"`
	} `json:"end_location" binding:"required"`

	DepartureTime   time.Time `json:"departure_time" binding:"required"`
	AvailableSeats  int       `json:"available_seats" binding:"required,min=1"`
	PricePerSeat    float64   `json:"price_per_seat" binding:"required,min=0"`
	AllowedDetourKm float64   `json:"allowed_detour_km" binding:"required,min=0"`
}

// CreateRideRequestRequest represents the request format for creating a ride request
type CreateRideRequestRequest struct {
	StartLocation struct {
		Latitude  float64 `json:"lat" binding:"required"`
		Longitude float64 `json:"lng" binding:"required"`
		Address   string  `json:"address" binding:"required"`
	} `json:"start_location" binding:"required"`

	EndLocation struct {
		Latitude  float64 `json:"lat" binding:"required"`
		Longitude float64 `json:"lng" binding:"required"`
		Address   string  `json:"address" binding:"required"`
	} `json:"end_location" binding:"required"`

	DepartureTime time.Time `json:"departure_time" binding:"required"`
	NumPassengers int       `json:"num_passengers" binding:"required,min=1"`
	MaxPrice      float64   `json:"max_price" binding:"required,min=0"`
}

// CreateRideOffer handles creating a new ride offer
func (h *RideHandler) CreateRideOffer(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var request CreateRideOfferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := h.rideService.CreateRideOffer(
		id,
		request.StartLocation.Latitude,
		request.StartLocation.Longitude,
		request.StartLocation.Address,
		request.EndLocation.Latitude,
		request.EndLocation.Longitude,
		request.EndLocation.Address,
		request.DepartureTime,
		request.AvailableSeats,
		request.PricePerSeat,
		request.AllowedDetourKm,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Ride offer created successfully",
		"ride_offer": offer,
	})
}

// CreateRideRequest handles creating a new ride request
func (h *RideHandler) CreateRideRequest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var request CreateRideRequestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rideRequest, err := h.rideService.CreateRideRequest(
		id,
		request.StartLocation.Latitude,
		request.StartLocation.Longitude,
		request.StartLocation.Address,
		request.EndLocation.Latitude,
		request.EndLocation.Longitude,
		request.EndLocation.Address,
		request.DepartureTime,
		request.NumPassengers,
		request.MaxPrice,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Ride request created successfully",
		"ride_request": rideRequest,
	})
}

// GetMyRideOffers handles retrieving all ride offers by the authenticated driver
func (h *RideHandler) GetMyRideOffers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	offers, err := h.rideService.GetRideOffersByDriver(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ride offers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ride_offers": offers,
	})
}

// GetMyRideRequests handles retrieving all ride requests by the authenticated passenger
func (h *RideHandler) GetMyRideRequests(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	requests, err := h.rideService.GetRideRequestsByPassenger(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ride requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ride_requests": requests,
	})
}

// ConfirmMatch handles confirming a ride match
func (h *RideHandler) ConfirmMatch(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	matchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match ID"})
		return
	}

	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	isDriver := roleStr == "driver" || roleStr == "both"

	if err := h.rideService.ConfirmMatch(matchID, id, isDriver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Match confirmed successfully",
	})
}
