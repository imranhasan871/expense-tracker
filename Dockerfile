# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

# Add ca-certificates for secure connections (RDS, etc.)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary and assets from builder
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
