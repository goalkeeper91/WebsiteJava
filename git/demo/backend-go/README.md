# Twitch Auth Backend (Go)

Clean Architecture Go-Backend mit Twitch OAuth2-Authentifizierung und Bot-Integration.

## Features

### ✅ Core Features (Implementiert)
- ✅ Twitch OAuth2 Login Flow
- ✅ User Management
- ✅ Session Management mit Cookies
- ✅ AES-GCM Verschlüsselung für Tokens
- ✅ Clean Architecture (Hexagonal)
- ✅ PostgreSQL mit Liquibase Migrations
- ✅ Docker & Docker Compose
- ✅ Nginx Reverse Proxy
# Twitch Auth Backend (Go)

Clean Architecture Go-Backend mit Twitch OAuth2-Authentifizierung und Bot-Integration.

## Features

### ✅ Core Features (Implementiert)
- ✅ Twitch OAuth2 Login Flow
- ✅ User Management
- ✅ Session Management mit Cookies
- ✅ AES-GCM Verschlüsselung für Tokens
- ✅ Clean Architecture (Hexagonal)
- ✅ PostgreSQL mit Liquibase Migrations
- ✅ Docker & Docker Compose
- ✅ Nginx Reverse Proxy

### 🚧 Extended Features (In Arbeit)
- 🔄 Redis Pub/Sub für Bot-Kommunikation
- 🔄 Chat Commands Management
- 🔄 Stream Management (Live-Status, Stats)
- 🔄 Activity Tracking mit WebSocket
- 🔄 Twitch API Integration
- 🔄 Contact Requests

## Projektstruktur

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Einstiegspunkt
├── internal/
│   ├── domain/
│   │   ├── user.go                 # User Entity
│   │   ├── session.go              # Session Entity
│   │   ├── auth_token.go           # Auth Token Entity
│   │   ├── chat_command.go         # Chat Command Entity
│   │   ├── stream_activity.go      # Stream Activity Entity
│   │   ├── twitch_channel.go       # Twitch Channel Entity
│   │   └── errors.go               # Domain Errors
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── auth_token_repository.go
│   │   ├── chat_command_repository.go
│   │   ├── stream_activity_repository.go
│   │   ├── twitch_channel_repository.go
│   │   └── postgres/
│   │       ├── user_repository.go
│   │       ├── auth_token_repository.go
│   │       └── chat_command_repository.go
│   ├── service/
│   │   └── auth_service.go
│   ├── handler/
│   │   └── auth_handler.go
│   ├── security/
│   │   └── crypto.go
│   └── infrastructure/
│       └── redis/
│           └── redis_service.go    # Redis Pub/Sub
├── pkg/
│   └── config/
│       └── config.go
├── migrations/
│   └── liquibase/
│       ├── changelog-master.xml
│       └── changesets/
├── docker/
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── go.mod
├── Makefile
├── README.md
├── SETUP.md
├── API.md
├── DEPLOYMENT.md
├── MIGRATION_GUIDE.md              # Java → Go Migration
└── STRUCTURE.md
```

## Tech Stack

- **Go 1.21+**: Backend Language
- **PostgreSQL 15**: Datenbank
- **Redis 7**: Pub/Sub & Caching
- **Liquibase**: Database Migrations
- **Nginx**: Reverse Proxy
- **Docker**: Containerization

## Dependencies

- `gorilla/mux` - HTTP Router
- `gorilla/sessions` - Session Management
- `gorilla/websocket` - WebSocket Support
- `redis/go-redis` - Redis Client
- `lib/pq` - PostgreSQL Driver
- `golang.org/x/oauth2` - OAuth2 Client

## Quick Start

### 1. Environment konfigurieren

```bash
cp .env.example .env
# .env bearbeiten und Werte eintragen
```

### 2. Mit Docker starten

```bash
docker-compose up -d
```

### 3. Lokal entwickeln

```bash
go mod download
go run cmd/server/main.go
```

## Environment Variables

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=twitch_auth
DB_USER=postgres
DB_PASSWORD=postgres

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Twitch OAuth
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
TWITCH_REDIRECT_URL=http://localhost/auth/callback

# Security
GO_APP_SECRET_KEY=your-32-character-secret-key!!
SESSION_SECRET=your-session-secret-key-min-32c

# Frontend
FRONTEND_URL=http://localhost:3000
```

## API Endpoints

### Auth
- `GET /auth/login` - Initiiert Twitch OAuth Flow
- `GET /auth/callback` - OAuth Callback Handler
- `POST /auth/logout` - Beendet Session
- `GET /auth/me` - Gibt aktuellen User zurück

### Dashboard (Coming Soon)
- `GET /api/dashboard/commands` - Chat Commands
- `GET /api/dashboard/stream/info` - Stream Infos
- `GET /api/dashboard/stream/activities` - Recent Activities

## Redis Pub/Sub

### Channels

- **bot:events** - Backend → Bot Signale
- **backend:events** - Bot → Backend Events

### Beispiel: Bot Signal senden

```go
redisService.SendJoinChannelSignal(twitchUserID)
redisService.SendRefreshCommandsSignal(twitchUserID)
```

## Make Commands

```bash
make help              # Zeigt alle Commands
make build             # Baut das Binary
make run               # Startet Server lokal
make test              # Führt Tests aus
make docker-up         # Startet Docker Container
make docker-down       # Stoppt Docker Container
make migrate-up        # Führt Migrations aus
make dev               # Startet Dev-Umgebung
```

## Migration von Java

Siehe [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) für Details zur Migration vom Java Spring Boot Backend.

## Dokumentation

- **[SETUP.md](./SETUP.md)** - Detaillierte Setup-Anleitung
- **[API.md](./API.md)** - API-Dokumentation
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Production Deployment
- **[MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)** - Java zu Go Migration
- **[STRUCTURE.md](./STRUCTURE.md)** - Architektur-Übersicht

## Entwicklung

### Tests ausführen

```bash
make test
```

### Code formatieren

```bash
make fmt
```

### Linter ausführen

```bash
make lint
```

## Lizenz

MIT
### 🚧 Extended Features (In Arbeit)
- 🔄 Redis Pub/Sub für Bot-Kommunikation
- 🔄 Chat Commands Management
- 🔄 Stream Management (Live-Status, Stats)
- 🔄 Activity Tracking mit WebSocket
- 🔄 Twitch API Integration
- 🔄 Contact Requests

## Projektstruktur

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Einstiegspunkt
├── internal/
│   ├── domain/
│   │   ├── user.go                 # User Entity
│   │   ├── session.go              # Session Entity
│   │   ├── auth_token.go           # Auth Token Entity
│   │   ├── chat_command.go         # Chat Command Entity
│   │   ├── stream_activity.go      # Stream Activity Entity
│   │   ├── twitch_channel.go       # Twitch Channel Entity
│   │   └── errors.go               # Domain Errors
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── auth_token_repository.go
│   │   ├── chat_command_repository.go
│   │   ├── stream_activity_repository.go
│   │   ├── twitch_channel_repository.go
│   │   └── postgres/
│   │       ├── user_repository.go
│   │       ├── auth_token_repository.go
│   │       └── chat_command_repository.go
│   ├── service/
│   │   └── auth_service.go
│   ├── handler/
│   │   └── auth_handler.go
│   ├── security/
│   │   └── crypto.go
│   └── infrastructure/
│       └── redis/
│           └── redis_service.go    # Redis Pub/Sub
├── pkg/
│   └── config/
│       └── config.go
├── migrations/
│   └── liquibase/
│       ├── changelog-master.xml
│       └── changesets/
├── docker/
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── go.mod
├── Makefile
├── README.md
├── SETUP.md
├── API.md
├── DEPLOYMENT.md
├── MIGRATION_GUIDE.md              # Java → Go Migration
└── STRUCTURE.md
```

## Tech Stack

- **Go 1.21+**: Backend Language
- **PostgreSQL 15**: Datenbank
- **Redis 7**: Pub/Sub & Caching
- **Liquibase**: Database Migrations
- **Nginx**: Reverse Proxy
- **Docker**: Containerization

## Dependencies

- `gorilla/mux` - HTTP Router
- `gorilla/sessions` - Session Management
- `gorilla/websocket` - WebSocket Support
- `redis/go-redis` - Redis Client
- `lib/pq` - PostgreSQL Driver
- `golang.org/x/oauth2` - OAuth2 Client

## Quick Start

### 1. Environment konfigurieren

```bash
cp .env.example .env
# .env bearbeiten und Werte eintragen
```

### 2. Mit Docker starten

```bash
docker-compose up -d
```

### 3. Lokal entwickeln

```bash
go mod download
go run cmd/server/main.go
```

## Environment Variables

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=twitch_auth
DB_USER=postgres
DB_PASSWORD=postgres

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Twitch OAuth
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
TWITCH_REDIRECT_URL=http://localhost/auth/callback

# Security
GO_APP_SECRET_KEY=your-32-character-secret-key!!
SESSION_SECRET=your-session-secret-key-min-32c

# Frontend
FRONTEND_URL=http://localhost:3000
```

## API Endpoints

### Auth
- `GET /auth/login` - Initiiert Twitch OAuth Flow
- `GET /auth/callback` - OAuth Callback Handler
- `POST /auth/logout` - Beendet Session
- `GET /auth/me` - Gibt aktuellen User zurück

### Dashboard (Coming Soon)
- `GET /api/dashboard/commands` - Chat Commands
- `GET /api/dashboard/stream/info` - Stream Infos
- `GET /api/dashboard/stream/activities` - Recent Activities

## Redis Pub/Sub

### Channels

- **bot:events** - Backend → Bot Signale
- **backend:events** - Bot → Backend Events

### Beispiel: Bot Signal senden

```go
redisService.SendJoinChannelSignal(twitchUserID)
redisService.SendRefreshCommandsSignal(twitchUserID)
```

## Make Commands

```bash
make help              # Zeigt alle Commands
make build             # Baut das Binary
make run               # Startet Server lokal
make test              # Führt Tests aus
make docker-up         # Startet Docker Container
make docker-down       # Stoppt Docker Container
make migrate-up        # Führt Migrations aus
make dev               # Startet Dev-Umgebung
```

## Migration von Java

Siehe [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) für Details zur Migration vom Java Spring Boot Backend.

## Dokumentation

- **[SETUP.md](./SETUP.md)** - Detaillierte Setup-Anleitung
- **[API.md](./API.md)** - API-Dokumentation
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Production Deployment
- **[MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)** - Java zu Go Migration
- **[STRUCTURE.md](./STRUCTURE.md)** - Architektur-Übersicht

## Entwicklung

### Tests ausführen

```bash
make test
```

### Code formatieren

```bash
make fmt
```

### Linter ausführen

```bash
make lint
```

## Lizenz

MIT