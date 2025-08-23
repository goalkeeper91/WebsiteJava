# === Build Stage ===
FROM maven:3.9-eclipse-temurin-24 AS build

WORKDIR /app

COPY . .

# Führen Sie den Maven-Build durch, um die ausführbare JAR-Datei zu erstellen
RUN mvn clean package -DskipTests

# === Runtime Stage ===
# Verwenden Sie ein schlankes JRE-Image für die Laufzeit
FROM eclipse-temurin:24-jre

# Fügen Sie Hosts-Einträge für die Ubuntu-Repositorys hinzu, um Verbindungsprobleme zu vermeiden
# Wir verwenden eine temporäre Datei innerhalb des RUN-Befehls.
RUN cat <<EOF > /tmp/hosts_additions && \
    cat /etc/hosts >> /tmp/hosts_additions && \
    mv /tmp/hosts_additions /etc/hosts
91.189.91.83 archive.ubuntu.com
91.189.91.81 security.ubuntu.com
EOF

# Installieren Sie den PostgreSQL-Client, da das Skript psql benötigt
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/target/*.jar app.jar

# Korrigierter Pfad: Da der Build-Kontext der Backend-Ordner ist, muss
# der Pfad nur der Dateiname sein.
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

EXPOSE 8080

# ENTRYPOINT verwendet jetzt das robuste Skript
ENTRYPOINT ["wait-for-it.sh", "postgres", "java", "-jar", "app.jar"]
