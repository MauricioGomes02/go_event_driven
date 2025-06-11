#!/bin/sh

set -e

# Wait for MySQL
echo "⏳ Waiting for MySQL at ${DATABASE_HOST}:${DATABASE_PORT}..."
until nc -z "$DATABASE_HOST" "$DATABASE_PORT"; do
  sleep 1
done
echo "✅ MySQL is available."

# Wait for Kafka (optional)
if [ -n "$KAFKA_BROKER" ]; then
  KAFKA_HOST=$(echo "$KAFKA_BROKER" | cut -d':' -f1)
  KAFKA_PORT=$(echo "$KAFKA_BROKER" | cut -d':' -f2)
  echo "⏳ Waiting for Kafka at ${KAFKA_HOST}:${KAFKA_PORT}..."
  until nc -z "$KAFKA_HOST" "$KAFKA_PORT"; do
    sleep 1
  done
  echo "✅ Kafka is available."
fi

# Execute the main process
exec "$@"
