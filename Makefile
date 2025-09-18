.PHONY: build run test clean docker-build docker-run migrate seed

APP_NAME=koperasi-merah-putih
BUILD_DIR=bin
MAIN_PATH=cmd/main.go

build:
	@echo "Building application..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	@echo "Running application..."
	@go run $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@golangci-lint run

migrate-up:
	@echo "Running database migrations..."
	@go run $(MAIN_PATH) -migrate

migrate-down:
	@echo "Rolling back database migrations..."
	@go run $(MAIN_PATH) -migrate-down

seed:
	@echo "Seeding database..."
	@go run $(MAIN_PATH) -seed

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(APP_NAME)

docker-compose-up:
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping services with docker-compose..."
	@docker-compose down

cassandra-init:
	@echo "Initializing Cassandra schema..."
	@cqlsh -f internal/models/cassandra/schemas.cql

postgres-init:
	@echo "Initializing PostgreSQL database..."
	@psql -h localhost -U postgres -d postgres -c "CREATE DATABASE IF NOT EXISTS koperasi_merah_putih;"

dev:
	@echo "Starting development environment..."
	@air

help:
	@echo "Available commands:"
	@echo "  build              Build the application"
	@echo "  run                Run the application"
	@echo "  test               Run tests"
	@echo "  clean              Clean build artifacts"
	@echo "  install-deps       Install Go dependencies"
	@echo "  lint               Run linter"
	@echo "  migrate-up         Run database migrations"
	@echo "  migrate-down       Rollback database migrations"
	@echo "  seed               Seed database with sample data"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run Docker container"
	@echo "  docker-compose-up  Start services with docker-compose"
	@echo "  docker-compose-down Stop services with docker-compose"
	@echo "  cassandra-init     Initialize Cassandra schema"
	@echo "  postgres-init      Initialize PostgreSQL database"
	@echo "  dev                Start development environment with hot reload"
	@echo "  help               Show this help message"