package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/ride-sharing-app/api/handlers"
	"github.com/yourusername/ride-sharing-app/api/middleware"
	"github.com/yourusername/ride-sharing-app/domain/model"
	"github.com/yourusername/ride-sharing-app/infrastructure/auth"
)

// Setup sets up the API routes
func Setup(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	rideHandler *handlers.RideHandler,
	jwtService *auth.JWTService,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Public routes
	router.POST("/api/v1/register", userHandler.Register)
	router.POST("/api/v1/login", userHandler.Login)

	// API v1 routes group
	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware(jwtService))
	{
		// User routes
		apiV1.GET("/profile", userHandler.GetProfile)

		// Driver routes
		driverRoutes := apiV1.Group("/driver")
		driverRoutes.Use(middleware.RoleMiddleware(model.RoleDriver, model.RoleBoth))
		{
			driverRoutes.POST("/profile", userHandler.RegisterDriverProfile)
			driverRoutes.GET("/profile", userHandler.GetDriverProfile)
			driverRoutes.POST("/rides", rideHandler.CreateRideOffer)
			driverRoutes.GET("/rides", rideHandler.GetMyRideOffers)
		}

		// Passenger routes
		passengerRoutes := apiV1.Group("/passenger")
		passengerRoutes.Use(middleware.RoleMiddleware(model.RolePassenger, model.RoleBoth))
		{
			passengerRoutes.POST("/rides", rideHandler.CreateRideRequest)
			passengerRoutes.GET("/rides", rideHandler.GetMyRideRequests)
		}

		// Match routes (available to both drivers and passengers)
		matchRoutes := apiV1.Group("/matches")
		{
			matchRoutes.POST("/:id/confirm", rideHandler.ConfirmMatch)
		}
	}
}
