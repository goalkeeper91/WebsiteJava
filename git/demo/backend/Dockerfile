# === Build Stage ===
FROM maven:3.9-eclipse-temurin-24 AS build

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

# === JRE-with-tools Stage ===
# Wir verwenden das vollständige JDK-Image, um jstack zu erhalten
FROM eclipse-temurin:24-jdk AS tools

# === Runtime Stage ===
FROM eclipse-temurin:24-jre

WORKDIR /app

COPY --from=build /app/target/*.jar app.jar

RUN apt-get update && apt-get install -y find && \
    JSTACK_PATH=$(find / -name "jstack" -print -quit) && \
    cp "$JSTACK_PATH" /usr/local/bin/jstack && \
    chmod +x /usr/local/bin/jstack && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

ENTRYPOINT ["/bin/bash", "-c", "until printf '' > /dev/tcp/postgres/5432 2>/dev/null; do echo 'Warte auf die PostgreSQL-Datenbank...'; sleep 1; done; exec java -jar app.jar"]
