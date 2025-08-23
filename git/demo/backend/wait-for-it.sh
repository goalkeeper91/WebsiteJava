#!/bin/bash
# wait-for-it.sh

# Warten, bis der PostgreSQL-Host bereit ist
# (hier verwenden wir 'postgres' als Host-Namen, wie im docker-compose.yml definiert)
host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  >&2 echo "Postgres ist noch nicht verfügbar - schläft"
  sleep 1
done

>&2 echo "Postgres ist bereit - Anwendung starten"
exec $cmd