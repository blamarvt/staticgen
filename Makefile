.PHONY: build test clean run all

# Binary name
BINARY_NAME=staticgen
OUTPUT_DIR=bin
CMD_DIR=cmd/staticgen

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean

# Build flags
LDFLAGS=-ldflags "-s -w"

all: clean test build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	$(GOTEST) -v ./tests/...
	@echo "Tests complete!"

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./tests/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(OUTPUT_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated!"

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Format complete!"

vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "Vet complete!"

lint: fmt vet
	@echo "Linting complete!"

help:
	@echo "Makefile commands:"
	@echo "  make build          - Build the binary"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make run            - Build and run the binary"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run fmt and vet"
	@echo "  make all            - Clean, test, and build"
