package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import postgres dialect
	"github.com/yourusername/ride-sharing-app/domain/model"
)

// InitDB initializes the database connection
func InitDB(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Enable Logger
	db.LogMode(true)

	// Auto-migrate the schema
	if err := migrateSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

// migrateSchema migrates the DB schema
func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.DriverProfile{},
		&model.RideOffer{},
		&model.RideRequest{},
		&model.RideMatch{},
	).Error
}
