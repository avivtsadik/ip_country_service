# IP Country Service

A production-grade REST API service that provides IP geolocation lookups with rate limiting and extensible datastore support.

## Features

- **REST API** with `/v1/find-country` endpoint
- **Custom rate limiting** using token bucket algorithm
- **Extensible datastore** interface (currently supports CSV)
- **Production-ready** with graceful shutdown and proper error handling
- **Comprehensive test suite** with unit and integration tests

## Setup

### 1. Environment Configuration

Create a `.env` file in the project root:

```env
PORT=8080
RATE_LIMIT_RPS=10.0
DATASTORE_TYPE=csv
DATASTORE_FILE=testdata/sample_ips.csv
```

**Environment Variables:**
- `PORT` - Server port (default: 8080)
- `RATE_LIMIT_RPS` - Requests per second limit (default: 10.0)
- `DATASTORE_TYPE` - Type of datastore (currently only "csv")
- `DATASTORE_FILE` - Path to CSV data file

### 2. Data File

Ensure your CSV file exists with the format:
```csv
8.8.8.8,Mountain View,United States
1.1.1.1,San Francisco,United States
```

The sample data file is provided at `testdata/sample_ips.csv`.

## Running the Service

### Start the Server

```bash
go run .
```

The service will start on the configured port and display:
```
IP Country Service starting on port 8080
Rate limit: 10.0 RPS
Datastore: csv (testdata/sample_ips.csv)
Server starting on port :8080
Endpoints available:
  GET /v1/find-country?ip=8.8.8.8
  GET /health
```

### API Usage

**Find country by IP:**
```bash
curl "http://localhost:8080/v1/find-country?ip=8.8.8.8"
```

**Response:**
```json
{
  "country": "United States",
  "city": "Mountain View"
}
```

**Health check:**
```bash
curl "http://localhost:8080/health"
```

**Response:**
```json
{
  "status": "ok"
}
```

## Running Tests

### All Tests
```bash
go test ./...
```

### Unit Tests Only
```bash
go test ./internal/...
```

### Integration Tests Only  
```bash
go test ./tests/...
```

### Verbose Output
```bash
go test ./... -v
```

### Specific Test Package
```bash
go test ./internal/services/ -v
go test ./internal/utils/ -v
go test ./tests/ -v
```

## Test Coverage

The project includes comprehensive tests:

- **Unit Tests** (16 tests) - Fast, isolated tests with mocks
  - Service layer business logic
  - IP validation utilities  
  - Safe goroutine wrapper
  - Rate limiter algorithm

- **Integration Tests** (14 tests) - End-to-end functionality
  - HTTP handlers with full request/response cycle
  - CSV datastore with file operations
  - Rate limiting middleware

## API Documentation

### `GET /v1/find-country`

**Query Parameters:**
- `ip` (required) - IPv4 address to lookup

**Success Response (200):**
```json
{
  "country": "United States",
  "city": "Mountain View"
}
```

**Error Responses:**
- `400 Bad Request` - Missing or invalid IP parameter
- `404 Not Found` - IP address not found in database
- `405 Method Not Allowed` - Only GET requests allowed
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### `GET /health`

**Success Response (200):**
```json
{
  "status": "ok"
}
```

## Architecture

```
├── main.go                 # Application entry point
├── internal/
│   ├── app/               # Application setup and configuration
│   ├── config/            # Environment variable configuration  
│   ├── datastores/        # Pluggable datastore implementations
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # Rate limiting middleware
│   ├── models/            # Data models
│   ├── services/          # Business logic layer
│   └── utils/             # Utility functions
├── tests/                 # Integration tests
├── testdata/              # Sample data files
└── .env                   # Environment configuration
```

## Development

The service follows clean architecture principles with clear separation of concerns:

- **Handlers** manage HTTP concerns
- **Services** contain business logic  
- **Datastores** provide pluggable data access
- **Middleware** handles cross-cutting concerns
- **Utils** provide reusable utilities

All components are unit tested with minimal dependencies for fast, reliable testing.