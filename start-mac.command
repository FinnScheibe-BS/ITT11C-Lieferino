#!/usr/bin/env bash
# 🚀 Startet das KOMPLETTE Lieferino-Projekt (Frontend + Backend + Datenbank) – macOS.
# Tipp: lässt sich per Doppelklick im Finder starten.
set -e
cd "$(dirname "$0")"

echo "🐳 Backend + Datenbank starten (Docker)..."
( cd backend && docker compose up -d --build )

echo "📦 Frontend-Abhängigkeiten installieren..."
npm install

echo "🌐 Der Browser öffnet sich gleich auf http://localhost:5173"
( sleep 8 && open http://localhost:5173 ) &

echo "🚀 Frontend starten... (zum Beenden: Strg+C)"
npm run dev
