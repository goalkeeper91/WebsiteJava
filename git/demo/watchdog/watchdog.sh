#!/bin/sh
# Uptime-Watchdog: prueft regelmaessig, ob nginx laeuft und die ueberwachten
# Domains erreichbar sind, und schickt bei Ausfall (und bei Wiederkehr) eine
# E-Mail ueber den bestehenden mailrelay-Container (Postfix -> Gmail).
#
# Zwei getrennte Pruefungen pro Domain, weil sie unterschiedliche Ausfaelle
# erkennen:
#
#  - "nginx": Anfrage direkt an den nginx-Container (Docker-DNS) mit dem
#    echten Hostnamen/SNI. Faengt ab: nginx gestoppt/abgestuerzt/Crash-Loop
#    (Container-Name nicht mehr aufloesbar bzw. keine Antwort), kaputte
#    Config, fehlendes/abgelaufenes Zertifikat - unabhaengig von DNS.
#  - "extern": Anfrage ueber die oeffentliche Domain (DNS -> Server-IP).
#    Faengt zusaetzlich ab: falscher/entfernter DNS-Eintrag, abgelaufene
#    Domain. Laeuft aber vom Server selbst aus - ein Totalausfall des
#    Servers/Netzes ist so NICHT erkennbar (dann kann dieser Container auch
#    keine Mail mehr schicken). Dafuer zusaetzlich einen externen Dienst
#    (z.B. UptimeRobot) einrichten.
#
# Ein Alarm geht erst nach FAIL_THRESHOLD aufeinanderfolgenden Fehlschlaegen
# raus (kein Alarm bei einem einzelnen Aussetzer, z.B. waehrend eines
# Deploys mit nginx --force-recreate). Solange der Ausfall anhaelt, kommt
# alle REMINDER_SECONDS eine Erinnerung; bei Wiederkehr eine Entwarnung.
# Als "gesendet" gilt ein Alarm erst, wenn der SMTP-Versand geklappt hat -
# schlaegt der Versand fehl, wird beim naechsten Durchlauf erneut versucht.

INTERVAL="${CHECK_INTERVAL:-60}"
FAIL_THRESHOLD="${FAIL_THRESHOLD:-2}"
REMINDER_SECONDS="${REMINDER_SECONDS:-21600}"
HOSTS="${WATCH_HOSTS:-fuchshilft.de}"
NGINX_HOST="${NGINX_HOST:-nginx}"
SMTP_URL="${SMTP_URL:-smtp://mailrelay:25}"
MAIL_TO="${ALERT_EMAIL_TO:-}"
MAIL_FROM="${ALERT_EMAIL_FROM:-watchdog@goalkeeper91.de}"
STATE_DIR="${STATE_DIR:-/tmp/watchdog-state}"

if [ -z "$MAIL_TO" ]; then
  echo "watchdog: ALERT_EMAIL_TO fehlt - Abbruch." >&2
  exit 1
fi

mkdir -p "$STATE_DIR"

log() {
  echo "[$(date -u +'%Y-%m-%d %H:%M:%S UTC')] $*"
}

# send_mail <betreff> <text> - Betreff bewusst nur ASCII (kein RFC-2047-
# Encoding noetig), Text als UTF-8. Zeilenenden als CRLF, wie SMTP es
# verlangt.
send_mail() {
  subject="$1"
  body="$2"
  mailfile="$(mktemp)"
  {
    printf 'From: Uptime-Watchdog <%s>\n' "$MAIL_FROM"
    printf 'To: %s\n' "$MAIL_TO"
    printf 'Subject: %s\n' "$subject"
    printf 'Date: %s\n' "$(date -u +'%a, %d %b %Y %H:%M:%S +0000')"
    printf 'MIME-Version: 1.0\n'
    printf 'Content-Type: text/plain; charset=UTF-8\n'
    printf 'Content-Transfer-Encoding: 8bit\n'
    printf '\n'
    printf '%s\n' "$body"
  } | sed 's/$/\r/' > "$mailfile"

  if curl -sS --max-time 30 --url "$SMTP_URL" \
        --mail-from "$MAIL_FROM" --mail-rcpt "$MAIL_TO" \
        -T "$mailfile" >/dev/null 2>"$mailfile.err"; then
    rm -f "$mailfile" "$mailfile.err"
    return 0
  fi
  log "ALARM-MAIL NICHT GESENDET ('$subject'): $(cat "$mailfile.err")"
  rm -f "$mailfile" "$mailfile.err"
  return 1
}

# check_nginx <host> -> gibt bei Fehler eine Beschreibung auf stdout aus.
check_nginx() {
  host="$1"
  ip="$(getent hosts "$NGINX_HOST" | awk '{print $1; exit}')"
  if [ -z "$ip" ]; then
    echo "Container '$NGINX_HOST' ist im Docker-Netz nicht aufloesbar (gestoppt, abgestuerzt oder in Crash-Loop)."
    return 1
  fi
  if err="$(curl -sS --fail --max-time 10 -o /dev/null \
        --resolve "$host:443:$ip" "https://$host/" 2>&1)"; then
    return 0
  fi
  echo "nginx ($ip) antwortet nicht korrekt auf https://$host/ : $err"
  return 1
}

# check_extern <host> -> gibt bei Fehler eine Beschreibung auf stdout aus.
check_extern() {
  host="$1"
  if err="$(curl -sS --fail --max-time 15 -o /dev/null "https://$host/" 2>&1)"; then
    return 0
  fi
  echo "https://$host/ ueber oeffentliches DNS nicht erreichbar: $err"
  return 1
}

# handle_result <id> <label> <rc> <detail>
handle_result() {
  id="$1"; label="$2"; rc="$3"; detail="$4"
  fails_file="$STATE_DIR/$id.fails"
  alerted_file="$STATE_DIR/$id.alerted"

  if [ "$rc" -eq 0 ]; then
    echo 0 > "$fails_file"
    if [ -f "$alerted_file" ]; then
      if send_mail "[OK] $label wieder erreichbar" \
"$label ist wieder erreichbar.

Ausfall begann ca.: $(date -u -d "@$(cat "$alerted_file")" +'%Y-%m-%d %H:%M UTC' 2>/dev/null || echo 'unbekannt')
Wieder erreichbar: $(date -u +'%Y-%m-%d %H:%M UTC')"; then
        rm -f "$alerted_file"
        log "OK: $label wieder erreichbar (Entwarnung gesendet)"
      fi
    fi
    return
  fi

  fails=$(( $(cat "$fails_file" 2>/dev/null || echo 0) + 1 ))
  echo "$fails" > "$fails_file"
  log "FEHLER ($fails/$FAIL_THRESHOLD): $label - $detail"

  [ "$fails" -ge "$FAIL_THRESHOLD" ] || return

  now="$(date +%s)"
  if [ ! -f "$alerted_file" ]; then
    if send_mail "[ALARM] $label nicht erreichbar" \
"$label ist nicht erreichbar (seit $fails Pruefungen in Folge, Intervall ${INTERVAL}s).

Details:
$detail

Zeitpunkt: $(date -u +'%Y-%m-%d %H:%M UTC')"; then
      echo "$now" > "$alerted_file"
      log "ALARM gesendet: $label"
    fi
  else
    last="$(cat "$alerted_file")"
    if [ $(( now - last )) -ge "$REMINDER_SECONDS" ]; then
      if send_mail "[ALARM-ERINNERUNG] $label weiterhin nicht erreichbar" \
"$label ist weiterhin nicht erreichbar.

Details:
$detail

Zeitpunkt: $(date -u +'%Y-%m-%d %H:%M UTC')"; then
        echo "$now" > "$alerted_file"
        log "Erinnerung gesendet: $label"
      fi
    fi
  fi
}

log "Watchdog gestartet: Hosts='$HOSTS', Intervall=${INTERVAL}s, Schwelle=$FAIL_THRESHOLD, Empfaenger=$MAIL_TO"

while true; do
  for host in $HOSTS; do
    detail="$(check_nginx "$host")"; rc=$?
    handle_result "nginx-$host" "nginx / $host (intern)" "$rc" "$detail"

    # Ist nginx selbst down, schlaegt die oeffentliche Pruefung zwangslaeufig
    # mit an - ein zweiter Alarm fuer denselben Vorfall waere nur Rauschen.
    # Sie wird deshalb erst wieder bewertet, wenn nginx intern antwortet
    # (dann zeigt ein Fehler tatsaechlich DNS/Domain/Netz an).
    if [ "$rc" -eq 0 ]; then
      detail="$(check_extern "$host")"; rc=$?
      handle_result "extern-$host" "$host (oeffentlich)" "$rc" "$detail"
    fi
  done
  sleep "$INTERVAL"
done
