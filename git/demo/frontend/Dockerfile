# === Build Stage ===
FROM node:22-alpine AS build

WORKDIR /app

# Nur package.json und package-lock.json kopieren, damit Docker-Caching funktioniert
COPY frontend/package*.json ./

RUN npm install

# Den Rest des Frontends kopieren
COPY frontend/. .

# Production Build erzeugen
RUN npm run build

# === Runtime Stage ===
FROM nginx:stable-alpine

# Production Build von React in Nginx Webroot kopieren
COPY --from=build /app/dist /usr/share/nginx/html

COPY ../nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Standardport für Nginx
EXPOSE 80

# Container starten Nginx im Vordergrund
CMD ["nginx", "-g", "daemon off;"]
