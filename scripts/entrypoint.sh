#!/bin/bash
set -e

# Run migrations if DB_URL is set
if [ -n "$DB_URL" ]; then
  echo "Running migrations..."
  /usr/local/bin/migrate -path /app/migrations -database "$DB_URL" up
fi

echo "Starting application..."
exec "$@"
