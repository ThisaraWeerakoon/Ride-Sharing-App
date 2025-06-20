FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ride-sharing-app .

# Use a minimal alpine image for the final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/ride-sharing-app .
COPY --from=builder /app/config/config.yaml ./config/

# Set environment variables
ENV GIN_MODE=release
ENV APP_SERVER_PORT=8080

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./ride-sharing-app"]
