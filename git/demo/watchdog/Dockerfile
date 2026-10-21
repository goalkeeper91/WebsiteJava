FROM alpine:3.20

# curl fuer die HTTP(S)-Pruefungen UND den SMTP-Versand an den mailrelay-
# Container (curl spricht SMTP), ca-certificates damit abgelaufene/ungueltige
# Zertifikate wirklich als Fehler erkannt werden. getent/awk/sed/date kommen
# aus Alpines Basis (musl-utils/busybox).
RUN apk add --no-cache curl ca-certificates

COPY watchdog.sh /watchdog.sh
RUN chmod +x /watchdog.sh

# Kein root noetig - der Container schreibt nur nach /tmp.
USER nobody

ENTRYPOINT ["/watchdog.sh"]
