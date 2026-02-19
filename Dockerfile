# ── Stage 1 : Build ──────────────────────────────────────────
FROM golang:1.18-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Copier les fichiers de dépendances en premier (cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler un binaire statique
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/arem-shop ./cmd/api

# ── Stage 2 : Run ───────────────────────────────────────────
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copier uniquement le binaire compilé
COPY --from=builder /bin/arem-shop /app/arem-shop

# Copier les migrations (utiles pour debug/init manuelle)
COPY --from=builder /src/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/arem-shop"]
