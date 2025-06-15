#!/bin/sh
set -e

echo "Waiting for PostgreSQL..."
until pg_isready -h postgres -U user -d coins; do sleep 1; done

echo "Applying migrations..."
for migration in /app/migrations/up/*.sql; do
  echo "Applying $migration"
  psql postgres://user:password@postgres:5432/coins?sslmode=disable -f "$migration"
done

exec "$@"