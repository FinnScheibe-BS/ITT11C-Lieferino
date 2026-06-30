@echo off
REM 🚀 Startet das KOMPLETTE Lieferino-Projekt (Frontend + Backend + Datenbank) - Windows.
cd /d "%~dp0"

echo 🐳 Backend + Datenbank starten (Docker)...
cd backend
docker compose up -d --build
cd ..

echo 📦 Frontend-Abhaengigkeiten installieren...
call npm install

echo 🌐 Der Browser oeffnet sich auf http://localhost:5173
start "" http://localhost:5173

echo 🚀 Frontend starten... (zum Beenden: Strg+C)
call npm run dev
