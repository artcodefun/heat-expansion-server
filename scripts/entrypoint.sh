#!/bin/bash
set -e

# Run migrations if GAME_DB_URL is set
if [ -n "$GAME_DB_URL" ]; then
  echo "Running migrations..."
  migrations_table="${GAME_MIGRATIONS_TABLE:-game_schema_migrations}"
  separator='&'
  if [[ "$GAME_DB_URL" != *"?"* ]]; then
    separator='?'
  fi
  /usr/local/bin/migrate -path /app/migrations -database "${GAME_DB_URL}${separator}x-migrations-table=${migrations_table}" up
fi

echo "Starting application..."
exec "$@"
