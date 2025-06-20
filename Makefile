.PHONY: build run clean docker-build docker-run docker-compose

# Build the application
build:
	go build -o bin/ride-sharing-app

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin

# Build Docker image
docker-build:
	docker build -t ride-sharing-app .

# Run Docker container
docker-run:
	docker run -p 8080:8080 ride-sharing-app

# Run with Docker Compose
docker-compose:
	docker-compose up -d

# Stop Docker Compose services
docker-compose-down:
	docker-compose down

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Install dependencies
deps:
	go mod tidy

# Generate Swagger documentation
swagger:
	swag init -g main.go -o api/docs
