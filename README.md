# cadprev_apis

Go backend application for consuming and analyzing APIs from https://apicadprev.trabalho.gov.br/api-docs/#

## Architecture

This project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** concepts:

### Layers

```
cmd/api/                    # Application entry point
└── main.go                 # Dependency injection and server setup

internal/
├── domain/                 # Domain layer (business logic, entities, interfaces)
│   ├── entity/            # Domain entities
│   │   └── data.go        # Data entity
│   ├── repository.go      # Repository interface
│   └── api_client.go      # External API client interface
│
├── usecase/               # Application layer (orchestration)
│   └── ingest.go          # Data ingestion use case
│
├── repository/            # Infrastructure layer (repository implementation)
│   └── inmemory.go        # In-memory stub repository
│
├── api/                   # Infrastructure layer (API client implementation)
│   └── stub_client.go     # Stub external API client
│
└── handler/               # Interface layer (HTTP handlers)
    └── data_handler.go    # HTTP handlers and routes
```

### Design Principles

1. **Domain Independence**: The domain layer has no dependencies on external frameworks or libraries
2. **Use Cases Orchestrate**: Business logic is orchestrated in the use case layer
3. **Interface-Based Design**: Repository and API clients are defined as interfaces in the domain layer
4. **Dependency Injection**: Dependencies are injected through constructors
5. **Context Propagation**: All operations use `context.Context` for cancellation and timeouts
6. **Standard Library**: Built with `net/http` and Go standard library

## Features

- Data ingestion from external API
- In-memory data persistence (stub implementation)
- RESTful API endpoints
- Graceful shutdown
- Health check endpoint
- Context-based cancellation
- Idiomatic Go code

## API Endpoints

### Health Check
```
GET /health
```
Returns server health status

### Ingest All Data
```
POST /api/ingest
```
Fetches all data from external API and persists it to the repository

### Ingest Specific Data
```
POST /api/ingest/{id}
```
Fetches specific data by ID from external API and persists it

### Get All Data
```
GET /api/data
```
Retrieves all data from the repository

### Get Data by ID
```
GET /api/data/{id}
```
Retrieves specific data by ID from the repository

## Getting Started

### Prerequisites

- Go 1.24+ installed

### Installation

```bash
# Clone the repository
git clone https://github.com/Filipefr15/cadprev_apis.git
cd cadprev_apis

# Download dependencies
go mod download

# Build the application
go build -o bin/api ./cmd/api
```

### Running the Application

```bash
# Run directly
go run cmd/api/main.go

# Or run the built binary
./bin/api
```

The server will start on port 8080 by default.

### Configuration

Environment variables:
- `PORT`: Server port (default: 8080)
- `EXTERNAL_API_URL`: External API base URL (default: https://apicadprev.trabalho.gov.br)

Example:
```bash
PORT=3000 EXTERNAL_API_URL=https://api.example.com go run cmd/api/main.go
```

## Usage Examples

### Check server health
```bash
curl http://localhost:8080/health
```

### Ingest all data
```bash
curl -X POST http://localhost:8080/api/ingest
```

### Ingest specific data
```bash
curl -X POST http://localhost:8080/api/ingest/custom-id-123
```

### Get all data
```bash
curl http://localhost:8080/api/data
```

### Get specific data
```bash
curl http://localhost:8080/api/data/ext-api-1
```

## Development

### Project Structure

The project follows the standard Go project layout with internal packages:

- `cmd/`: Application entry points
- `internal/`: Private application code
  - `domain/`: Domain models and interfaces (core business logic)
  - `usecase/`: Application use cases (orchestration)
  - `repository/`: Data persistence implementations
  - `api/`: External API client implementations
  - `handler/`: HTTP handlers and routing

### Adding New Features

1. **New Domain Entity**: Add to `internal/domain/entity/`
2. **New Interface**: Define in `internal/domain/`
3. **New Use Case**: Implement in `internal/usecase/`
4. **New Repository**: Implement in `internal/repository/`
5. **New Handler**: Add to `internal/handler/` and register routes

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run

# Vet code
go vet ./...

# Build and verify compilation
go build ./...
```

## Current Implementation Status

⚠️ **Note**: This is a stub implementation for demonstration purposes:

- **API Client**: Currently returns mock data instead of making real HTTP calls
- **Repository**: Uses in-memory storage instead of a real database

### Production Readiness

To make this production-ready:

1. Replace `StubAPIClient` with a real HTTP client using `net/http`
2. Replace `InMemoryRepository` with a real database (PostgreSQL, MongoDB, etc.)
3. Add proper error handling and logging
4. Add authentication and authorization
5. Add rate limiting and circuit breakers
6. Add monitoring and observability
7. Add comprehensive unit and integration tests
8. Add API documentation (OpenAPI/Swagger)

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

## License

This project follows the license terms of the original repository.

## Contributing

Contributions are welcome! Please follow Go best practices and maintain the Clean Architecture principles.
