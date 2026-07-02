# ☸️ Lieferino auf Kubernetes (mit Patroni-HA-Datenbank)

Diese Anleitung bringt Lieferino **skalierbar** auf einen Kubernetes-Cluster:
- die **PostgreSQL-Datenbank** läuft als **Patroni-Cluster** (3 Knoten, automatischer Failover)
- das **Backend** läuft in mehreren Kopien hinter einem Service

> Für die normale lokale Entwicklung reicht weiterhin `docker compose` (siehe `../SETUP.md`).
> Kubernetes ist die „große", produktionsnahe Variante.

---

## 0. Voraussetzungen
- ein Kubernetes-Cluster – lokal z. B. **kind** oder **minikube**
- **kubectl** (mit dem Cluster verbunden)
- **Docker** (zum Bauen des Backend-Images)

---

## 1. Backend-Image bauen + in den Cluster laden
Das Backend-Image wird lokal gebaut (enthält alle nötigen Assets):
```sh
docker build -t lieferino-backend:latest ./backend
```
Damit der Cluster das lokale Image nutzen kann:
```sh
# kind:
kind load docker-image lieferino-backend:latest

# minikube:
minikube image load lieferino-backend:latest
```

## 2. Patroni-Operator installieren (einmalig)
Wir nutzen den bewährten **Zalando postgres-operator**. Der baut aus unserer
kleinen Beschreibung automatisch einen vollständigen Patroni-HA-Cluster:
```sh
kubectl apply -k github.com/zalando/postgres-operator/manifests
```
Kurz warten, bis der Operator läuft:
```sh
kubectl wait --for=condition=available --timeout=120s deploy/postgres-operator
```

## 3. Lieferino ausrollen
```sh
kubectl apply -f k8s/00-namespace.yaml
kubectl apply -f k8s/10-postgres.yaml     # Patroni-DB-Cluster (3 Knoten)
kubectl apply -f k8s/20-backend.yaml      # Backend (2 Kopien)
```
Auf den DB-Cluster warten (dauert beim ersten Mal ein paar Minuten):
```sh
kubectl -n lieferino get postgresql
kubectl -n lieferino get pods -w
```

## 4. App erreichen
Backend per Port-Forward auf den lokalen Port 8090 holen:
```sh
kubectl -n lieferino port-forward svc/lieferino-backend 8090:8090
```
Test: <http://localhost:8090/health> → `{"status":"ok"}`.

Das **Frontend** läuft wie gewohnt lokal (`npm run dev`) und spricht mit
`http://localhost:8090` – die erlaubte Herkunft `http://localhost:5173` ist im
Backend bereits gesetzt (`CORS_ORIGINS`).

---

## 5. 🧪 Failover testen (das Herzstück von Patroni)
Aktuellen Leader + Replicas anzeigen:
```sh
kubectl -n lieferino get pods -l application=spilo -L spilo-role
```
Den **Leader** absichtlich löschen:
```sh
kubectl -n lieferino delete pod <name-des-leaders>
```
Innerhalb weniger Sekunden macht Patroni eine Replica zum neuen Leader, und der
Service `lieferino-db` zeigt automatisch auf den neuen Leader. Die App läuft
weiter – **ohne dass Daten verloren gehen**.

## 6. 📈 Skalieren
```sh
# Mehr Backend-Kopien:
kubectl -n lieferino scale deploy/lieferino-backend --replicas=4

# Mehr DB-Knoten (in k8s/10-postgres.yaml "numberOfInstances" erhöhen, dann):
kubectl apply -f k8s/10-postgres.yaml
```

## 7. Aufräumen
```sh
kubectl delete -f k8s/20-backend.yaml -f k8s/10-postgres.yaml -f k8s/00-namespace.yaml
```

---

## ℹ️ Hinweise
- Den DB-Service `lieferino-db` (Leader, Schreiben) bzw. `lieferino-db-repl`
  (Replicas, Lesen) erzeugt der Operator automatisch.
- Das Passwort des App-Nutzers `lieferino` generiert der Operator in ein Secret;
  das Backend liest es direkt von dort (siehe `20-backend.yaml`).
- Die Manifeste hier sind bewusst klein gehalten und auf lokale Cluster
  (kind/minikube) abgestimmt (geringe CPU/RAM-Anforderungen).
