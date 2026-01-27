# Go Pastebin Example

A simple pastebin web application written in Go with unit and integration tests.

## Features

- Create text pastes with optional expiration
- View and delete pastes
- RESTful API
- Minimal frontend interface
- In-memory storage
- Comprehensive test coverage

## Security Notice

⚠️ **This repository contains a vulnerable dependency for demonstration purposes.**

This project uses `gin-gonic/gin v1.6.3` which contains **CVE-2020-28483** (directory traversal vulnerability). This is intentional for testing and educational purposes. Do not use this in production.

## Prerequisites

- Go 1.20 or higher

## Installation

```bash
# Clone the repository
git clone https://github.com/colin-harness/go-example.git
cd go-example

# Install dependencies
go mod download
```

## Running the Application

```bash
# Start the server (default port 8080)
go run main.go

# Or specify a custom port
PORT=3000 go run main.go
```

Visit `http://localhost:8080` in your browser to use the pastebin.

## API Endpoints

- `POST /paste` - Create a new paste
  ```json
  {
    "content": "your text here",
    "ttl": 3600  // optional, seconds until expiration
  }
  ```

- `GET /api/paste/:id` - Retrieve a paste (JSON)
- `GET /paste/:id` - View a paste (HTML)
- `DELETE /api/paste/:id` - Delete a paste

## Running Tests

```bash
# Run all tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./store
go test -v ./handler
```

The test suite includes:
- Unit tests for the storage layer (100% coverage)
- Unit tests for HTTP handlers
- Integration tests for full workflows
- Concurrency tests
- Tests take >10 seconds to ensure thorough validation

## Project Structure

```
.
├── main.go              # Application entry point
├── handler/             # HTTP handlers
│   ├── handler.go
│   └── handler_test.go
├── store/               # Data storage layer
│   ├── store.go
│   └── store_test.go
├── static/              # Frontend assets
│   └── index.html
├── integration_test.go  # End-to-end tests
├── go.mod
└── README.md
```

## Coverage

The project maintains high test coverage:
- Store package: 100%
- Handler package: ~57%
- Overall integration: ~47%

## Development

```bash
# Run tests in watch mode (requires entr)
find . -name '*.go' | entr -c go test ./...

# Format code
go fmt ./...

# Run static analysis
go vet ./...
```

## Known Issues

- Uses vulnerable dependency gin v1.6.3 (CVE-2020-28483)
- In-memory storage only (data lost on restart)
- No authentication or authorization
- No rate limiting

## License

MIT
