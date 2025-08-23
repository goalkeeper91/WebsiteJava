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

COPY --from=tools /opt/java/openjdk/bin/jstack /usr/local/bin/jstack
RUN chmod +x /usr/local/bin/jstack

EXPOSE 8080

ENTRYPOINT ["/bin/bash", "-c", "until printf '' > /dev/tcp/postgres/5432 2>/dev/null; do echo 'Warte auf die PostgreSQL-Datenbank...'; sleep 1; done; exec java -jar app.jar"]
