# Go Koperasi Management System
# Makefile for development tasks

.PHONY: help build run test migrate seed fresh-migrate dev-setup clean

help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  migrate       - Run database migrations"
	@echo "  seed          - Run database seeders"
	@echo "  migrate-fresh - Drop all tables, migrate, and seed"
	@echo "  dev-setup     - Setup development environment"
	@echo "  clean         - Clean build artifacts"

build:
	@echo "Building application..."
	go build -o bin/main cmd/main.go

run:
	@echo "Running application..."
	go run cmd/main.go

test:
	@echo "Running tests..."
	go test ./... -v

migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

seed:
	@echo "Running database seeders..."
	go run cmd/seeder/main.go

migrate-fresh:
	@echo "Running fresh migration..."
	go run cmd/migrate/main.go -fresh

migrate-drop:
	@echo "Dropping tables and migrating..."
	go run cmd/migrate/main.go -drop

dev-setup:
	@echo "Setting up development environment..."
	@echo "Installing dependencies..."
	go mod download
	@echo "Running fresh migration..."
	go run cmd/migrate/main.go -fresh
	@echo "Development setup complete!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

fmt:
	@echo "Formatting code..."
	go fmt ./...