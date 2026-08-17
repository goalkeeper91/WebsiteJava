FROM debian:bookworm-slim

# Nur Postfix als reiner Weiterleitungs-Relay (siehe main.cf) - kein
# Dovecot/IMAP, keine echten Postfaecher, deshalb bewusst schlank.
RUN apt-get update && apt-get install -y --no-install-recommends \
      postfix libsasl2-modules ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY main.cf /etc/postfix/main.cf
COPY virtual /etc/postfix/virtual
COPY entrypoint.sh /entrypoint.sh

# virtual-Map ist nicht geheim (nur E-Mail-Adressen, keine Passwoerter) und
# kann deshalb schon beim Build gepostmapt werden. sasl_passwd (das
# Gmail-App-Passwort) ist geheim und wird bewusst erst beim Container-Start
# aus Umgebungsvariablen erzeugt (siehe entrypoint.sh) statt beim Build -
# damit landet das Passwort nie im Docker-Image/Layer-Cache oder in git.
RUN chmod +x /entrypoint.sh \
    && postmap /etc/postfix/virtual \
    # Chroot fuer alle Postfix-Dienste aus (siehe Debugging: der chrootete
    # smtp-Ausgangsprozess sieht /etc/resolv.conf nicht, DNS-Aufloesung fuer
    # smtp.gmail.com schlug deshalb mit "Host not found" fehl, obwohl der
    # Name ganz normal aufloesbar ist). Der Container selbst ist hier schon
    # die Isolation, ein zusaetzliches Chroot bringt in diesem engen
    # Einzelzweck-Setup keinen relevanten Zusatznutzen.
    && postconf -F '*/*/chroot=n'

EXPOSE 25
ENTRYPOINT ["/entrypoint.sh"]
