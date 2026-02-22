include .env
export

# Settings
BINARY_NAME=weatherbot
GOOSE_BIN=$(shell go env GOPATH)/bin/goose

.PHONY: all help fmt clean test run build up down status create \
        docker-up docker-down docker-logs docker-db docker-clean install-goose

all: help

## --- COMMON ---

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -hE '^##' $(MAKEFILE_LIST) | sed -e 's/## //' | awk -F ':' '{printf "  %-15s %s\n", $$1, $$2}'

## fmt: Format source code according to Go standards
fmt:
	go fmt ./...

## test: Run all project tests with verbose output
test:
	go test -v ./...

## clean: Remove compiled binary and temporary files
clean:
	rm -f $(BINARY_NAME)

## --- DEVELOPMENT (Local) ---

## build: Build the binary locally
build:
	go build -o $(BINARY_NAME) main.go

## run: Build and run the application locally
run: build
	./$(BINARY_NAME)

## up: Run database migrations locally
up: install-goose
	$(GOOSE_BIN) up

## down: Rollback the last local migration
down: install-goose
	$(GOOSE_BIN) down

## status: Show local migration status
status: install-goose
	$(GOOSE_BIN) status

## create: Create a new migration (usage: make create name=migration_name)
create: install-goose
	@if [ -z "$(name)" ]; then echo "Error: name is required. Use 'make create name=your_migration_name'"; exit 1; fi
	$(GOOSE_BIN) create $(name) sql

## install-goose: Install goose migration tool if not present
install-goose:
	@ls $(GOOSE_BIN) > /dev/null 2>&1 || (echo "Goose not found. Installing..." && go install github.com/pressly/goose/v3/cmd/goose@latest)

## --- DOCKER (Production-like) ---

## docker-up: Build and start the entire stack in Docker (Bot + DB)
docker-up:
	docker-compose up --build -d

## docker-down: Stop and remove all project containers
docker-down:
	docker-compose down

## docker-logs: Follow logs from all containers
docker-logs:
	docker-compose logs -f

## docker-db: Enter PostgreSQL shell inside the Docker container
docker-db:
	docker exec -it weather-db psql -U postgres

## docker-clean: Deep clean of unused Docker resources (prune)
docker-clean:
	docker system prune -f
