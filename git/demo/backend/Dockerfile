# === Build Stage ===
FROM maven:3.9-eclipse-temurin-21 AS build
# Hinweis: Temurin 24 ist sehr neu (Early Access), 21 ist der aktuelle LTS-Standard.

WORKDIR /app

# 1. Nur die POMs kopieren, um Layer-Caching für Dependencies zu nutzen
COPY pom.xml .
COPY backend/pom.xml backend/

# 2. Nur den tatsächlichen Source-Code des Backends kopieren
# Das verhindert, dass "Geisterdateien" aus dem Root oder anderen Ordnern mitrutschen
COPY backend/src backend/src

# 3. Den Build explizit für das Backend ausführen
# Wir gehen in den Backend-Ordner, damit Maven lokal dort baut
WORKDIR /app/backend
RUN mvn clean package -DskipTests

# === Runtime Stage ===
FROM eclipse-temurin:21-jre
WORKDIR /app

# Der Pfad zur JAR hat sich durch das 'WORKDIR' oben geändert:
COPY --from=build /app/backend/target/*.jar app.jar

# Korrigierter Pfad: Da der Build-Kontext der Backend-Ordner ist, muss
# der Pfad nur der Dateiname sein.
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

EXPOSE 8080

# ENTRYPOINT verwendet jetzt das robuste Skript
ENTRYPOINT ["wait-for-it.sh", "postgres", "java", "-jar", "app.jar"]
