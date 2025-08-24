#!/bin/bash
# wait-for-it.sh

# Variablen für den Host-Namen und den auszuführenden Befehl
host="$1"
shift
cmd="$@"

# Warte auf das Vorhandensein der Datenbank und der notwendigen Tabelle
# (Verwendet psql um auf die Verfügbarkeit der 'twitch_commands' Tabelle zu warten)
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  >&2 echo "Postgres ist noch nicht verfügbar - schläft"
  sleep 1
done

>&2 echo "Datenbank und Schema sind bereit. Anwendung starten."
exec $cmd
