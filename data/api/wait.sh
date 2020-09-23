#!/bin/sh
# wait-for-postgres.sh

set -e

cmd="$@"
  
until PGPASSWORD=$POSTGRES_PASSWORD psql -h $DB_CONTAINER -U $POSTGRES_USER -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done
  
>&2 echo "Postgres is up - executing command"
exec $cmd
