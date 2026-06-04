# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git curl

# Install golang-migrate — separate layer, cached until version changes
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/heat-expansion-server ./cmd/server

# Final stage
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates and bash (useful for entrypoint scripts)
RUN apk add --no-cache ca-certificates bash

# Copy binary and migrate tool
COPY --from=builder /app/heat-expansion-server /app/
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy migrations
COPY internal/game/infrastructure/db/migrations /app/migrations/game
COPY internal/auth/infrastructure/db/migrations /app/migrations/auth
COPY internal/billing/infrastructure/db/migrations /app/migrations/billing
COPY internal/admin/infrastructure/db/migrations /app/migrations/admin

# Copy entrypoint script
COPY scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expose ports
EXPOSE 8080 8081 8082 8083

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./heat-expansion-server"]
