#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# scripts/docker-start.sh — Démarrage complet via Docker
# ──────────────────────────────────────────────────────────────
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

log() { printf '\033[1;34m[docker]\033[0m %s\n' "$*"; }
fail() { printf '\033[1;31m[docker][error]\033[0m %s\n' "$*" >&2; exit 1; }

# ── Vérifier Docker ───────────────────────────────────────────
command -v docker >/dev/null 2>&1 || fail "Docker n'est pas installé"
docker info >/dev/null 2>&1 || fail "Le démon Docker n'est pas lancé"

# ── Détecter docker compose (v2 prioritaire sur v1) ───────────
if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
else
  fail "docker compose n'est pas installé"
fi

log "Utilisation de : $COMPOSE_CMD"

# ── Créer .env si absent ──────────────────────────────────────
if [[ ! -f "$ROOT_DIR/.env" ]]; then
  cp "$ROOT_DIR/.env.example" "$ROOT_DIR/.env"
  log ".env créé depuis .env.example"

  # Auto-générer un JWT_SECRET fort
  if command -v openssl >/dev/null 2>&1; then
    generated_secret="$(openssl rand -hex 32)"
    sed -i "s|^JWT_SECRET=.*|JWT_SECRET=${generated_secret}|" "$ROOT_DIR/.env"
    log "JWT_SECRET auto-généré"
  fi
else
  log ".env existant détecté"
fi

# ── Build & Start ─────────────────────────────────────────────
log "Build et démarrage des conteneurs..."
$COMPOSE_CMD up --build "$@"
