# Development Dockerfile for Go backend
FROM golang:alpine AS base

# Install air for hot reload in development
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Use air for hot reload in development
CMD ["air", "-c", ".air.toml"] 