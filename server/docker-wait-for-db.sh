#!/bin/bash
set -e

host="$1"
shift
port="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASS psql -h "$host" -U "$DB_USER" -d "$PGDATABASE" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec env "$@"
