#!/bin/sh
set -e

# Das Gmail-App-Passwort kommt ausschliesslich als Umgebungsvariable rein
# (aus der .env auf dem Server, NICHT aus git - siehe README/Dockerfile-
# Kommentar) und wird erst hier zur Laufzeit in die von Postfix erwartete
# sasl_passwd-Datei geschrieben.
if [ -z "$GMAIL_RELAY_USER" ] || [ -z "$GMAIL_RELAY_APP_PASSWORD" ]; then
  echo "entrypoint.sh: GMAIL_RELAY_USER / GMAIL_RELAY_APP_PASSWORD fehlen - Abbruch." >&2
  exit 1
fi

echo "[smtp.gmail.com]:587 ${GMAIL_RELAY_USER}:${GMAIL_RELAY_APP_PASSWORD}" > /etc/postfix/sasl_passwd
postmap /etc/postfix/sasl_passwd
chmod 600 /etc/postfix/sasl_passwd /etc/postfix/sasl_passwd.db
chown root:root /etc/postfix/sasl_passwd /etc/postfix/sasl_passwd.db

exec postfix start-fg
