.PHONY: build run dev-up dev-down clean test

APP_NAME := family-svc
BUILD_DIR := ./bin

## build: Build the Go binary
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

## run: Run the service locally
run:
	go run ./cmd/server

## test: Run all project tests
test:
	go test -v ./...

## dev-up: Start local dev dependencies (PostgreSQL)
dev-up:
	docker compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec postgres pg_isready -U postgres -d familydb > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL is ready!"

## dev-down: Stop local dev dependencies
dev-down:
	docker compose down

## clean: Remove build artifacts and volumes
clean:
	rm -rf $(BUILD_DIR)
	docker compose down -v
