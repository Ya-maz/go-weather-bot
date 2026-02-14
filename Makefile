include .env
export

# Binary name
BINARY_NAME=weatherbot
GOOSE_BIN=$(shell go env GOPATH)/bin/goose

.PHONY: all build run clean help up down status create install-goose

all: help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -hE '^##' $(MAKEFILE_LIST) | sed -e 's/## //' | awk -F ':' '{printf "  %-15s %s\n", $$1, $$2}'

## install-goose: Install goose migration tool if not present
install-goose:
	@ls $(GOOSE_BIN) > /dev/null 2>&1 || (echo "Goose not found. Installing..." && go install github.com/pressly/goose/v3/cmd/goose@latest)

## build: Build the binary
build:
	go build -o $(BINARY_NAME) main.go

## run: Build and run the application
run: build
	./$(BINARY_NAME)

## up: Run database migrations
up: install-goose
	$(GOOSE_BIN) up

## down: Rollback the last migration
down: install-goose
	$(GOOSE_BIN) down

## status: Show migration status
status: install-goose
	$(GOOSE_BIN) status

## create: Create a new migration (usage: make create name=migration_name)
create: install-goose
	@if [ -z "$(name)" ]; then echo "Error: name is required. Use 'make create name=your_migration_name'"; exit 1; fi
	$(GOOSE_BIN) create $(name) sql

## clean: Remove binary
clean:
	rm -f $(BINARY_NAME)
