# === Build Stage ===
FROM maven:3.9-eclipse-temurin-24 AS build

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

# === Runtime Stage ===
FROM eclipse-temurin:24-jdk

WORKDIR /app

COPY --from=build /app/target/*.jar app.jar

# Fügen Sie das wait-for-it-Skript hinzu und machen Sie es ausführbar
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

EXPOSE 8080

CMD ["java", "-jar", "app.jar"]
