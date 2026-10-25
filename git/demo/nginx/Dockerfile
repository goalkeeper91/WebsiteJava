# Wraps the stock nginx image with a periodic self-reload, so renewed
# certificates (written to disk by the separate certbot renewer container,
# which has no hook back into this container) actually get picked up.
# Without this, nginx keeps serving whatever certificate was in memory at
# last start/reload indefinitely, silently going stale/expired once a
# renewal happens on disk - this is what caused goalkeeper91.de and
# timer.goalkeeper91.de to serve an expired certificate despite valid files
# already sitting under ./certs.
# 1.25 ist eine abgelaufene Mainline-Linie ohne Sicherheitskorrekturen; 1.30 ist die aktuelle stabile Linie (Stand 25.10.2026).
FROM nginx:1.30-alpine

# Version nicht im Server-Header und in Fehlerseiten nennen (conf.d wird im http-Kontext eingelesen).
COPY 00-security.conf /etc/nginx/conf.d/00-security.conf

# Reload twice a day (matches the certbot renewer's 12h `certbot renew`
# loop) - cheap/idempotent when nothing actually changed, bounds staleness
# to at most ~12h after any real renewal.
RUN echo "0 */12 * * * nginx -s reload" > /etc/crontabs/root

CMD ["sh", "-c", "crond && nginx -g 'daemon off;'"]
