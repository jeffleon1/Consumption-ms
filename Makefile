.PHONY: all build run clean

APP_NAME := consumption-ms
SRC_DIR := ./cmd/api

run:
	@echo "Running $(APP_NAME)..."
	@go run $(SRC_DIR)/main.go

run-docker:
	@echo "Running $(APP_NAME) in Docker..."
	@docker compose up -d --build

down-docker:
	@echo "Stopping $(APP_NAME)..."
	@docker compose down