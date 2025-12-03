# Conslabs 2025
# Makefile automated build and run

include .env

# Variables
NAME = go-backend-nba
BUILD_PATH = ./build
MAIN_DIR = ./cmd/
MAIN_FILE = main.go
SEEDER_FILE=cmd/seeder/main.go
MIGRATION_PATH = ./db/migrations

# Target: build
.PHONY: build
build:
	@echo "🚀 Building the application..."
	@go build -o $(BUILD_PATH)/$(NAME) $(MAIN_DIR)*go
	@echo "Build completed. Output at $(BUILD_PATH)/ with name '$(NAME)'"

# Target: run without build
.PHONY: run
run:
	@echo "🚀 Running the application without building..."
	go run $(MAIN_DIR)$(MAIN_FILE)
	@echo "Run completed."

# Target: run prebuild program
.PHONY: start
start:
	@echo "🚀 Running the application from build file..."
	@cd $(BUILD_PATH) && ./$(NAME)

.PHONY: build-run
build-run:
	@echo "🚀 Building the application..."
	@go build -o $(BUILD_PATH)/$(NAME) $(MAIN_DIR)$(MAIN_FILE)
	@echo "Build completed. Output at $(BUILD_PATH)/$(NAME) with name $(NAME)"
	@echo "Starting application"
	@cd $(BUILD_PATH) && ./$(NAME)
	@echo "Build and run completed"	

# Start DB
.PHONY: run-db
run-db:
	@echo "🚀 Starting PostgreSQL container in db/"
	@cd db && sudo docker compose up -d

.PHONY: stop-db
stop-db:
	@echo "🚀 Stopping PostgreSQL container in db/"
	@cd db && sudo docker compose down

# Migrations
.PHONY: migrate-create
migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-fix
migrate-fix:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-seed
migrate-seed: ## Run the seeder script
	@echo "🌱 Starting seeder..."
	@go run $(SEEDER_FILE) --gseed

# This avoids "No rule to make target" error for extra args
%:
	@: