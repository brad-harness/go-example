.PHONY: build test test-coverage run clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove built binaries and coverage files"

# Build the application
build:
	go build -o pastebin .

# Run the application
run:
	go run main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -f pastebin coverage.out coverage.html
	go clean
