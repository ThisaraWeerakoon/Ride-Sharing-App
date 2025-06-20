package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ride-sharing-app/domain/model"
	"github.com/yourusername/ride-sharing-app/infrastructure/auth"
	"github.com/yourusername/ride-sharing-app/service"
)

// UserHandler handles user-related API requests
type UserHandler struct {
	userService *service.UserService
	jwtService  *auth.JWTService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService, jwtService *auth.JWTService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// RegisterRequest represents the request format for user registration
type RegisterRequest struct {
	FirstName string         `json:"first_name" binding:"required"`
	LastName  string         `json:"last_name" binding:"required"`
	Email     string         `json:"email" binding:"required,email"`
	Password  string         `json:"password" binding:"required,min=6"`
	Phone     string         `json:"phone" binding:"required"`
	Role      model.UserRole `json:"role" binding:"required,oneof=passenger driver both"`
}

// LoginRequest represents the request format for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterDriverRequest represents the request format for driver registration
type RegisterDriverRequest struct {
	LicenseNo  string `json:"license_no" binding:"required"`
	CarModel   string `json:"car_model" binding:"required"`
	CarPlateNo string `json:"car_plate_no" binding:"required"`
	NumSeats   int    `json:"num_seats" binding:"required,min=2"`
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
		request.Phone,
		request.Role,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
		"token":   token,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": user.ID,
		"token":   token,
		"role":    user.Role,
	})
}

// GetProfile handles retrieving user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
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

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
	})
}

// RegisterDriverProfile handles driver profile registration
func (h *UserHandler) RegisterDriverProfile(c *gin.Context) {
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

	var request RegisterDriverRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.userService.RegisterDriverProfile(
		id,
		request.LicenseNo,
		request.CarModel,
		request.CarPlateNo,
		request.NumSeats,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Driver profile registered successfully",
		"profile": profile,
	})
}

// GetDriverProfile handles retrieving driver profile
func (h *UserHandler) GetDriverProfile(c *gin.Context) {
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

	profile, err := h.userService.GetDriverProfile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get driver profile"})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
