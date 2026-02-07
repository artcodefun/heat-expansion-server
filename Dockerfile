# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git curl

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /app/heat-expansion-api ./cmd/api

# Install golang-migrate for the final image
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and bash (useful for entrypoint scripts)
RUN apk add --no-cache ca-certificates bash

# Copy binary and migrate tool
COPY --from=builder /app/heat-expansion-api /app/
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy migrations
COPY internal/infrastructure/db/migrations /app/migrations

# Copy entrypoint script
COPY scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expose port
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./heat-expansion-api"]
