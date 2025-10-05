#!/bin/bash
# Wait for PostgreSQL to be ready before running migrations

set -e

host="$1"
user="$2"
password="$3"
database="$4"
shift 4
cmd="$@"

echo "Waiting for PostgreSQL at $host to be ready..."

export PGPASSWORD="$password"

until psql -h "$host" -U "$user" -d "$database" -c '\q' 2>/dev/null; do
  >&2 echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

>&2 echo "PostgreSQL is up - executing command"
exec $cmd
