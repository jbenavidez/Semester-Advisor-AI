APP_NAME=semester-advisor-ai-service
APP_CMD=./cmd/web

.PHONY: up up_build down build logs clean restart shell lint lint_fix

up:
	@echo "Starting Semester Advisor AI services..."
	APP_CMD=${APP_CMD} docker compose up
	@echo "Semester Advisor AI services started!"

up_build:
	@echo "Stopping Semester Advisor AI services if running..."
	docker compose down
	@echo "Building and starting Semester Advisor AI services..."
	APP_CMD=${APP_CMD} docker compose up --build
	@echo "Semester Advisor AI services started!"

build:
	@echo "Building Semester Advisor AI Docker images..."
	APP_CMD=${APP_CMD} docker compose build
	@echo "Build completed!"

down:
	@echo "Stopping Semester Advisor AI services..."
	docker compose down

restart: down up_build

logs:
	docker compose logs -f

shell:
	docker exec -it ${APP_NAME} sh

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

lint_fix:
	@echo "Running golangci-lint with automatic fixes..."
	golangci-lint run --fix ./...

clean:
	@echo "Removing Semester Advisor AI containers and volumes..."
	docker compose down -v
	@echo "Cleanup completed!"