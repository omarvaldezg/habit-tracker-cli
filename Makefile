.PHONY: build clean install test run

BINARY_NAME=habit-tracker
VERSION=1.0.0
BUILD_DIR=build
GO_BUILD_FLAGS=-ldflags="-X main.version=$(VERSION)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/habit-tracker

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Install the application to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(GO_BUILD_FLAGS) ./cmd/habit-tracker

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run ./cmd/habit-tracker/main.go
