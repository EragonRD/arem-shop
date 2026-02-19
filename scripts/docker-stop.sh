#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# scripts/docker-stop.sh — Arrêt des conteneurs
# ──────────────────────────────────────────────────────────────
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

log() { printf '\033[1;34m[docker]\033[0m %s\n' "$*"; }

# ── Détecter docker compose (v2 prioritaire sur v1) ───────────
if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
else
  echo "docker compose introuvable" >&2; exit 1
fi

REMOVE_VOLUMES=0
if [[ "${1:-}" == "--volumes" ]] || [[ "${1:-}" == "-v" ]]; then
  REMOVE_VOLUMES=1
fi

if [[ "$REMOVE_VOLUMES" == "1" ]]; then
  log "Arrêt des conteneurs + suppression des volumes..."
  $COMPOSE_CMD down -v
else
  log "Arrêt des conteneurs (données DB conservées)..."
  $COMPOSE_CMD down
fi

log "Terminé ✅"
