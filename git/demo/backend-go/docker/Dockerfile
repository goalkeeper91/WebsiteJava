# Build Stage
FROM golang:1.21-alpine AS builder

# Build-Dependencies installieren
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Dependencies cachen
COPY go.mod go.sum ./
RUN go mod download

# Source Code kopieren
COPY . .

# Binary bauen
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/server \
    ./cmd/server

# Runtime Stage
FROM alpine:latest

# Runtime-Dependencies
RUN apk --no-cache add ca-certificates tzdata

# Non-root User erstellen
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Binary vom Builder kopieren
COPY --from=builder /build/server /app/server

# Ownership ändern
RUN chown -R appuser:appuser /app

# Als non-root User ausführen
USER appuser

# Port exponieren
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Server starten
CMD ["/app/server"]