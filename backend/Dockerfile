# --- Build-Stufe: Go-Programm im Container kompilieren ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Quellcode kopieren (inkl. der eingebetteten data/seed.bin).
COPY . .

# Eingebettete Assets prüfen.
RUN echo "d7a4819acd094cdf01827cd7e08bd868c322e95f7b399b3f1aa61a8f8fefc6de  data/seed.bin" | sha256sum -c -

# Abhängigkeiten auflösen + Programm bauen.
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o server .

# --- Lauf-Stufe: schlankes Image nur mit dem fertigen Programm ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
