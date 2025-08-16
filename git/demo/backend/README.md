# Streamer Website Backend

## Übersicht

Dieses Backend dient als zentraler Server für die **Streamer-Website**.
Es integriert sowohl **Discord- als auch Twitch-Bots**, verwaltet deren Status, Befehle und Konfigurationen und stellt REST-Endpoints für das Frontend bereit.

Die Architektur folgt einer **Service-Controller-Struktur**, die sauberes Refactoring und einfache Erweiterbarkeit ermöglicht.

---

## Package-Struktur

### `client`
- Enthält die Haupt-Klassen für die Bots (z. B. `DiscordBot`).
- Startet die Bots, verbindet sich mit den jeweiligen APIs und registriert Event-Handler.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/client)

### `commands/discord` & `commands/twitch`
- Definieren die verfügbaren **Commands** für die Bots.
- Jeder Command implementiert ein Interface (`command`) mit `getName()` und `execute()`.
- Neue Commands können einfach hinzugefügt werden.
- [Discord](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/commands/discord)
- [Twitch](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/commands/twitch)

### `config`
- Enthält Spring-spezifische Konfigurationen (z. B. Security, Application Properties).
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/config)

### `controller/discord` & `controller/twitch`
- REST-Controller für die Bots.
- Beispiel: `DiscordStatusController` liefert Statusinformationen wie „läuft seit“, aktive Channels etc.
- Schnittstelle für das Frontend, um Status abzufragen.
- [Discord](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/controller/discord)
- [Twitch](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/controller/twitch)

### `controller`
- Allgemeine Controller, die nicht direkt Bot-spezifisch sind.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/controller)

### `dto`
- Data Transfer Objects für REST-Kommunikation.
- Standardisiert die Daten zwischen Backend und Frontend.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/dto)

### `entity`
- Datenbank-Entities.
- Repräsentieren persistente Daten wie User, Tokens, Bot-Settings.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/entity)

### `repository`
- Spring Data JPA Repositories für CRUD-Operationen auf Entities.
- Ermöglicht einfache Datenbankinteraktion.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/repository)

### `security`
- Authentifizierung und Autorisierung für das Backend.
- JWT oder andere Security-Mechanismen.
- [Link](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/security)

### `service/discord` & `service/twitch`
- Logik für Discord- bzw. Twitch-Bot: Command Handling, Status-Verwaltung etc.
- `CommandService`: zentrale Dispatch-Logik für Commands.
- `StatusService`: hält den aktuellen Bot-Status (running, uptime, channels).
- [Discord](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/service/discord)
- [Twitch](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/service/twitch)
- [Service Root](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo/service)

### Root Packages
- [Demo Root](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/java/streamer_website/demo)
- [Resources](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend/src/main/resources)
- [Backend Root](https://github.com/goalkeeper91/WebsiteJava/tree/main/git/demo/backend)

---

## Beispiel-Flow: Discord-Bot

1. `DiscordBot` wird bei App-Start als Spring-Bean instanziiert.
2. Token wird aus `application.properties` gelesen.
3. Bot loggt sich ein (`DiscordClient.login().block()`).
4. `CommandService` registriert alle Commands und reagiert auf `MessageCreateEvent`.
5. `StatusService` hält den Status für REST-Controller verfügbar.
6. Frontend fragt `/api/discord/status` ab, um Laufzeit und aktive Channels anzuzeigen.

---

## Setup

1. **Dependencies installieren**:

```bash
./mvnw clean install
```

## Setup mit Docker

1. **Docker Image bauen**:

```bash
docker compose up --build -d
```
