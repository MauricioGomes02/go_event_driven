#!/bin/sh

set -e

# Wait for MySQL
echo "⏳ Waiting for MySQL at ${DATABASE_HOST}:${DATABASE_PORT}..."
until nc -z "$DATABASE_HOST" "$DATABASE_PORT"; do
  sleep 1
done
echo "✅ MySQL is available."

# Execute the main process
exec "$@"
