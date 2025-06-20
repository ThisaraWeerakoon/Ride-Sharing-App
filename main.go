package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ride-sharing-app/api/handlers"
	"github.com/yourusername/ride-sharing-app/api/routes"
	"github.com/yourusername/ride-sharing-app/config"
	"github.com/yourusername/ride-sharing-app/infrastructure/auth"
	"github.com/yourusername/ride-sharing-app/infrastructure/database"
	"github.com/yourusername/ride-sharing-app/repository"
	"github.com/yourusername/ride-sharing-app/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create repositories
	userRepo := repository.NewGormUserRepository(db)
	rideRepo := repository.NewGormRideRepository(db)

	// Create services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)
	userService := service.NewUserService(userRepo)
	rideService := service.NewRideService(rideRepo, userRepo)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService, jwtService)
	rideHandler := handlers.NewRideHandler(rideService)

	// Initialize Gin
	router := gin.Default()

	// Setup routes
	routes.Setup(router, userHandler, rideHandler, jwtService)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
