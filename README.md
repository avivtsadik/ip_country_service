# IP Country Service

A production-grade REST API service that provides IP geolocation lookups with rate limiting and extensible datastore support.

## Features

- **REST API** with `/v1/find-country` endpoint
- **Custom rate limiting** using token bucket algorithm
- **Extensible datastore** interface (supports CSV and JSON formats)
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
- `DATASTORE_TYPE` - Type of datastore ("csv" or "json", default: "csv")
- `DATASTORE_FILE` - Path to data file (CSV or JSON format)

**IDE Configuration (GoLand/IntelliJ):**
1. Create `.env` file with your configuration
2. Open **Run/Debug Configurations** (top-right dropdown next to ▶️)
3. Choose **Edit Configurations...**
4. Load environment variables from your `.env` file
5. Alternatively, set environment variables manually in the IDE

**Note:** JSON datastore support was added as a bonus feature to demonstrate the service's extensibility to multiple data sources.

### 2. Data File

**CSV Format:**
```csv
8.8.8.8,Mountain View,United States
1.1.1.1,San Francisco,United States
```

**JSON Format (Bonus Feature):**
```json
[
  {"ip": "8.8.4.4", "city": "Mountain View", "country": "United States"},
  {"ip": "1.0.0.1", "city": "Research", "country": "Australia"}
]
```

Sample data files are provided:
- `testdata/sample_ips.csv` (default)
- `testdata/sample_ips.json` (bonus extensibility demo)

**Switching Between Datastores:**
```bash
# Use CSV datastore
DATASTORE_TYPE=csv DATASTORE_FILE=testdata/sample_ips.csv go run .

# Use JSON datastore (bonus feature)
DATASTORE_TYPE=json DATASTORE_FILE=testdata/sample_ips.json go run .
```

## Running the Service

### Option 1: Direct Go Run

```bash
go run .
```

### Option 2: Docker Container

**Build and run with scripts:**
```bash
./build.sh    # Build Docker image
./run.sh      # Run container
```

**Manual Docker commands:**
```bash
# Build image
docker build -t ip-country-service .

# Run container
docker run -p 8080:8080 ip-country-service
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
go test ./internal/app/ -v
```

### Verbose Output
```bash
go test ./... -v
```

### Specific Test Package
```bash
go test ./internal/services/ -v
go test ./internal/utils/ -v
go test ./internal/app/ -v
go test ./internal/datastores/ -v
go test ./internal/middleware/ -v
```

## Test Coverage

The project includes comprehensive tests:

- **Unit Tests** (27 tests) - Fast, isolated tests with mocks
  - Service layer business logic
  - IP validation utilities  
  - Safe goroutine wrapper
  - Rate limiter algorithm
  - CSV datastore functionality
  - JSON datastore functionality (bonus)

- **Integration Tests** (7 tests) - End-to-end functionality
  - HTTP handlers with full request/response cycle
  - Rate limiting middleware
  - Full application flow testing

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
├── Dockerfile              # Docker container definition
├── build.sh                # Docker build script
├── run.sh                  # Docker run script
├── internal/
│   ├── app/               # Application setup and integration tests
│   ├── config/            # Environment variable configuration  
│   ├── datastores/        # Pluggable datastore implementations
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # Rate limiting middleware
│   ├── models/            # Data models
│   ├── services/          # Business logic layer
│   └── utils/             # Utility functions
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

## Bonus Features

### 1. **Comprehensive Testing** ✅
- **Unit Tests**: 27 tests covering all components with mocks
- **Integration Tests**: 7 end-to-end tests with full HTTP request/response cycle
- **Test Organization**: Tests follow Go conventions (next to source code)
- **Safe Goroutines**: All test goroutines use panic recovery
- **Context Support**: Full context propagation throughout the application

### 2. **Docker Containerization** ✅
- Multi-stage Dockerfile with optimized build using Go 1.24-alpine
- Executable build and run scripts (`build.sh`, `run.sh`)
- Production-ready container with proper environment configuration

### 3. **Multiple Datastore Support** ✅ (Extensibility Demo)
- **JSON Datastore**: Added as bonus to demonstrate extensible architecture
- **Environment-driven switching**: Simple `DATASTORE_TYPE` configuration
- **Identical interface**: Both CSV and JSON implement the same `DataStore` interface
- **Comprehensive testing**: 8 additional tests for JSON datastore implementation
- **Different sample data**: Each datastore has unique test data to demonstrate switching

This demonstrates how easily new data sources can be added (databases, APIs, etc.) without changing the core application logic.