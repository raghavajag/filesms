# Variables
APP_NAME=file-sharing-platform
BUILD_DIR=./build
MAIN_FILE=./cmd/api/main.go

# Targets
.PHONY: all build run clean test

all: clean build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build completed."

run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_FILE)

migrate:
	@echo "Running database migrations..."
	@go run ./db/migrate.go
	@echo "Migrations completed."
