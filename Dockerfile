# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy testdata directory
COPY --from=builder /app/testdata ./testdata

# Expose port
EXPOSE 8080

# Set environment variables with defaults
ENV PORT=8080
ENV RATE_LIMIT_RPS=10.0
ENV DATASTORE_TYPE=csv
ENV DATASTORE_FILE=testdata/sample_ips.csv

# Run the binary
CMD ["./main"]