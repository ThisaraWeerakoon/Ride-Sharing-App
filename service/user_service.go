package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/ride-sharing-app/domain/model"
	"github.com/yourusername/ride-sharing-app/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser registers a new user in the system
func (s *UserService) RegisterUser(firstName, lastName, email, password, phone string, role model.UserRole) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
		Phone:     phone,
		Role:      role,
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user by email and password
func (s *UserService) AuthenticateUser(email, password string) (*model.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(id uuid.UUID, firstName, lastName, phone string) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone = phone

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// RegisterDriverProfile registers a driver profile for a user
func (s *UserService) RegisterDriverProfile(
	userID uuid.UUID,
	licenseNo,
	carModel,
	carPlateNo string,
	numSeats int,
) (*model.DriverProfile, error) {
	// Check if user exists and can be a driver
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Role != model.RoleDriver && user.Role != model.RoleBoth {
		return nil, errors.New("user is not registered as a driver")
	}

	// Check if driver profile already exists
	existingProfile, err := s.userRepo.GetDriverProfile(userID)
	if err != nil {
		return nil, err
	}
	if existingProfile != nil {
		return nil, errors.New("driver profile already exists for this user")
	}

	// Create new driver profile
	profile := &model.DriverProfile{
		UserID:     userID,
		LicenseNo:  licenseNo,
		CarModel:   carModel,
		CarPlateNo: carPlateNo,
		NumSeats:   numSeats,
	}

	// Save profile to database
	if err := s.userRepo.CreateDriverProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// GetDriverProfile retrieves a driver's profile
func (s *UserService) GetDriverProfile(userID uuid.UUID) (*model.DriverProfile, error) {
	return s.userRepo.GetDriverProfile(userID)
}

// UpdateDriverProfile updates a driver's profile
func (s *UserService) UpdateDriverProfile(
	userID uuid.UUID,
	licenseNo,
	carModel,
	carPlateNo string,
	numSeats int,
) (*model.DriverProfile, error) {
	profile, err := s.userRepo.GetDriverProfile(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("driver profile not found")
	}

	profile.LicenseNo = licenseNo
	profile.CarModel = carModel
	profile.CarPlateNo = carPlateNo
	profile.NumSeats = numSeats

	if err := s.userRepo.UpdateDriverProfile(profile); err != nil {
		return nil, err
	}

	return profile, nil
}
