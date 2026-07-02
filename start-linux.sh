#!/usr/bin/env bash
# 🚀 Startet das KOMPLETTE Lieferino-Projekt (Frontend + Backend + Datenbank) – Linux.
set -e
cd "$(dirname "$0")"

echo "🐳 Backend + Datenbank starten (Docker)..."
( cd backend && docker compose up -d --build )

echo "📦 Frontend-Abhängigkeiten installieren..."
npm install

echo "🌐 Der Browser öffnet sich gleich auf http://localhost:5173"
( sleep 8 && xdg-open http://localhost:5173 >/dev/null 2>&1 ) &

echo "🚀 Frontend starten... (zum Beenden: Strg+C)"
npm run dev
