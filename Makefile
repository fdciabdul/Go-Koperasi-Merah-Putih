# Koperasi Merah Putih Development Commands

.PHONY: help dev build run test clean deps migrate

help:
	@echo "Available commands:"
	@echo "  make dev     - Start development server with hot reload"
	@echo "  make build   - Build the application"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make deps    - Install dependencies"
	@echo "  make migrate - Run database migrations"

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

dev:
	@echo "Starting development server with hot reload..."
	air

build:
	@echo "Building application..."
	go build -o bin/koperasi-app cmd/main.go

run:
	@echo "Starting server..."
	go run cmd/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf tmp/
	rm -rf bin/
	rm -f build-errors.log

migrate:
	@echo "Running migrations..."
	go run cmd/migrate/main.go
