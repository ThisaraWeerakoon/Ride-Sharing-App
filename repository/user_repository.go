package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/yourusername/ride-sharing-app/domain/model"
	repo "github.com/yourusername/ride-sharing-app/domain/repository"
)

// GormUserRepository is an implementation of UserRepository using Gorm
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB) repo.UserRepository {
	return &GormUserRepository{db: db}
}

// Create adds a new user to the database
func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves a user by ID
func (r *GormUserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *GormUserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *GormUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user from the database
func (r *GormUserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

// CreateDriverProfile adds a new driver profile to the database
func (r *GormUserRepository) CreateDriverProfile(profile *model.DriverProfile) error {
	return r.db.Create(profile).Error
}

// GetDriverProfile retrieves a driver profile by user ID
func (r *GormUserRepository) GetDriverProfile(userID uuid.UUID) (*model.DriverProfile, error) {
	var profile model.DriverProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// UpdateDriverProfile updates a driver profile in the database
func (r *GormUserRepository) UpdateDriverProfile(profile *model.DriverProfile) error {
	return r.db.Save(profile).Error
}
