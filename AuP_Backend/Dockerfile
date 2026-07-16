FROM golang:1.26-alpine

WORKDIR /app

# Modul-Dateien kopieren
COPY go.mod ./
# Da wir noch keine go.sum haben, lassen wir Go sie gleich selbst erzeugen

COPY . .

# Externe Pakete im Container laden
RUN go mod tidy

# Kompilieren für Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

EXPOSE 8080
CMD ["/server"]
