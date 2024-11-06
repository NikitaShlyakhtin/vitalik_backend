# Simple Makefile for a Go project

# Build the application
all: build test

.PHONY: build
build:
	@echo 'Building cmd/api...'

	go build -o=./bin/vitalik ./cmd/vitalik

	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/vitalik ./cmd/vitalik

.PHONY: run
run:
	@go run ./cmd/vitalik

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch
