# === Build Stage ===
FROM node:22-alpine AS build

WORKDIR /app

# Nur package.json und package-lock.json kopieren (Caching für npm install)
COPY frontend/package*.json ./

RUN npm install

# Restlichen Code kopieren
COPY frontend/. .

# Production Build erzeugen
RUN npm run build

