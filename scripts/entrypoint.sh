#!/bin/bash
set -e

# Run migrations if GAME_DB_URL is set
if [ -n "$GAME_DB_URL" ]; then
  echo "Running game migrations..."
  migrations_table="${GAME_MIGRATIONS_TABLE:-game_schema_migrations}"
  separator='&'
  if [[ "$GAME_DB_URL" != *"?"* ]]; then
    separator='?'
  fi
  /usr/local/bin/migrate -path /app/migrations/game -database "${GAME_DB_URL}${separator}x-migrations-table=${migrations_table}" up
fi

# Run auth migrations if AUTH_DB_URL is set
if [ -n "$AUTH_DB_URL" ]; then
  echo "Running auth migrations..."
  migrations_table="${AUTH_MIGRATIONS_TABLE:-auth_schema_migrations}"
  separator='&'
  if [[ "$AUTH_DB_URL" != *"?"* ]]; then
    separator='?'
  fi
  /usr/local/bin/migrate -path /app/migrations/auth -database "${AUTH_DB_URL}${separator}x-migrations-table=${migrations_table}" up
fi

# Run billing migrations if BILLING_DB_URL is set
if [ -n "$BILLING_DB_URL" ]; then
  echo "Running billing migrations..."
  migrations_table="${BILLING_MIGRATIONS_TABLE:-billing_schema_migrations}"
  separator='&'
  if [[ "$BILLING_DB_URL" != *"?"* ]]; then
    separator='?'
  fi
  /usr/local/bin/migrate -path /app/migrations/billing -database "${BILLING_DB_URL}${separator}x-migrations-table=${migrations_table}" up
fi

# Run admin migrations if ADMIN_DB_URL is set
if [ -n "$ADMIN_DB_URL" ]; then
  echo "Running admin migrations..."
  migrations_table="${ADMIN_MIGRATIONS_TABLE:-admin_schema_migrations}"
  separator='&'
  if [[ "$ADMIN_DB_URL" != *"?"* ]]; then
    separator='?'
  fi
  /usr/local/bin/migrate -path /app/migrations/admin -database "${ADMIN_DB_URL}${separator}x-migrations-table=${migrations_table}" up
fi

echo "Starting application..."
exec "$@"
