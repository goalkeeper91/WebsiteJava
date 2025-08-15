DemoJavaWebApp
📖 Projektbeschreibung

DemoJavaWebApp ist ein Showcase-Projekt, das meinen aktuellen Fortschritt in Java Backend-Entwicklung demonstriert. Das Projekt nutzt React und Tailwind im Frontend und ist vollständig mit Docker containerisiert, um die Anwendung einfach zu starten und zu verwalten.

Dieses Projekt dient als Beispiel für potenzielle Arbeitgeber, um meine Fähigkeiten in modernen Webtechnologien und Java zu sehen.

🏗️ Projektstruktur
demo/
├── backend/       # Java Backend
├── frontend/      # React Frontend mit Tailwind
└── docker-compose.yml  # Startet Backend und Frontend zusammen

⚙️ Technologien

Java – Backend-Logik

Spring Boot / Servlets – REST-API (Backend)

React – Frontend

Tailwind CSS – Styling

Docker & Docker Compose – Containerisierung und einfacher Start

Maven – Backend-Build-Tool

npm / yarn – Frontend-Build-Tool

🚀 Installation & Start
Voraussetzungen

Docker & Docker Compose

Node.js / npm (optional, falls Frontend lokal gebaut wird)

Projekt starten

Im Projektverzeichnis (demo) einfach:

docker-compose up --build


Backend läuft auf http://localhost:8080

Frontend läuft auf http://localhost:3000 (oder wie in docker-compose.yml konfiguriert)

📝 Nutzung

Nach dem Start sind Backend-API und Frontend direkt verfügbar.

Das Frontend konsumiert die Backend-API automatisch.

📂 Wichtige Dateien

backend/pom.xml – Maven-Konfiguration

frontend/package.json – Frontend-Abhängigkeiten & Scripts

docker-compose.yml – Container-Setup

📄 Lizenz

Copyright (c) 2025 [Dein Name].  
All rights reserved.  

Dieses Projekt darf nicht kopiert, verbreitet oder verändert werden ohne ausdrückliche schriftliche Genehmigung des Autors.
