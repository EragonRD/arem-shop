#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# scripts/docker-logs.sh — Afficher les logs des conteneurs
# ──────────────────────────────────────────────────────────────
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# ── Détecter docker compose (v2 prioritaire sur v1) ───────────
if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD="docker-compose"
else
  echo "docker compose introuvable" >&2; exit 1
fi

# Par défaut : suivre tous les services. Sinon le service passé en argument.
SERVICE="${1:-}"
$COMPOSE_CMD logs -f $SERVICE
