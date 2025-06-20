# Ride-Sharing-App

A ride-sharing application built with Go and Gin framework, with a monolithic architecture designed to scale to microservices in the future.

## Architecture Overview

The application uses a layered architecture with the following components:

1. **API Layer**: Handles HTTP requests/responses using Gin framework
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Manages data access
4. **Domain Layer**: Defines core business models
5. **Infrastructure**: Cross-cutting concerns like logging, auth, and configuration

## Current Monolithic Structure

```
ride-sharing-app/
├── api/            # API handlers and routes
├── config/         # Application configuration
├── domain/         # Business models and interfaces
├── infrastructure/ # Cross-cutting concerns
├── repository/     # Data access implementations
├── service/        # Business logic implementations
├── main.go         # Application entry point
└── go.mod          # Go module definition
```

## Future Microservices Path

The application is designed to be separated into the following potential microservices:

1. **User Service**: User registration, authentication, and profile management
2. **Ride Service**: Ride creation, matching, and management
3. **Notification Service**: User notifications
4. **Payment Service**: Payment processing
5. **Analytics Service**: Usage statistics and reporting

## Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL (or your preferred database)

### Installation

1. Clone the repository

   ```
   git clone https://github.com/yourusername/ride-sharing-app.git
   ```

2. Install dependencies

   ```
   go mod download
   ```

3. Configure environment variables (see config/config.go)

4. Run the application
   ```
   go run main.go
   ```

## API Documentation

API documentation is available at `/swagger/index.html` when the server is running.
