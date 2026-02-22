#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'USAGE'
Usage:
  scripts/first_run.sh [--no-run]

Options:
  --no-run   Prepare environment and database without starting the API.

Optional environment variables for initial seed (if no SuperAdmin exists):
  FIRST_SHOP_ID
  FIRST_SHOP_NAME
  FIRST_SHOP_WHATSAPP
  FIRST_SUPERADMIN_NAME
  FIRST_SUPERADMIN_EMAIL
  FIRST_SUPERADMIN_PASSWORD
USAGE
}

log() {
  printf '[first-run] %s\n' "$*"
}

fail() {
  printf '[first-run][error] %s\n' "$*" >&2
  exit 1
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    fail "Command not found: $cmd"
  fi
}

generate_uuid() {
  if command -v uuidgen >/dev/null 2>&1; then
    uuidgen | tr '[:upper:]' '[:lower:]'
    return
  fi

  if [[ -r /proc/sys/kernel/random/uuid ]]; then
    cat /proc/sys/kernel/random/uuid
    return
  fi

  fail "Unable to generate UUID (uuidgen not found and /proc/sys/kernel/random/uuid unavailable)"
}

NO_RUN=0
if [[ "${1:-}" == "--no-run" ]]; then
  NO_RUN=1
  shift
fi

if [[ $# -gt 0 ]]; then
  usage
  fail "Unknown argument(s): $*"
fi

require_cmd go
require_cmd psql
require_cmd createdb

if [[ ! -f "$ROOT_DIR/.env" ]]; then
  cp "$ROOT_DIR/config/.env.example" "$ROOT_DIR/.env"
  log ".env created from config/.env.example"

  # Auto-generate a stronger JWT secret for first run.
  if command -v openssl >/dev/null 2>&1; then
    generated_secret="$(openssl rand -hex 32)"
    sed -i "s|^JWT_SECRET=.*|JWT_SECRET=${generated_secret}|" "$ROOT_DIR/.env"
    log "JWT_SECRET auto-generated in .env"
  fi
fi

set -a
# shellcheck disable=SC1091
source "$ROOT_DIR/.env"
set +a

: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_USER:=postgres}"
: "${DB_NAME:=arem_shop}"
: "${DB_SSLMODE:=disable}"
: "${BCRYPT_COST:=12}"

export PGPASSWORD="${DB_PASSWORD:-}"

PG_BASE_CONN="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} sslmode=${DB_SSLMODE}"
POSTGRES_DB_CONN="${PG_BASE_CONN} dbname=postgres"
APP_DB_CONN="${PG_BASE_CONN} dbname=${DB_NAME}"

log "Checking PostgreSQL connectivity..."
psql "$POSTGRES_DB_CONN" -tAc "SELECT 1" >/dev/null

log "Ensuring database '${DB_NAME}' exists..."
DB_EXISTS="$(psql "$POSTGRES_DB_CONN" -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | tr -d '[:space:]')"
if [[ "$DB_EXISTS" != "1" ]]; then
  createdb --host "$DB_HOST" --port "$DB_PORT" --username "$DB_USER" "$DB_NAME"
  log "Database created: ${DB_NAME}"
else
  log "Database already exists: ${DB_NAME}"
fi

MIGRATION_FILE="$ROOT_DIR/migrations/000001_init.up.sql"
[[ -f "$MIGRATION_FILE" ]] || fail "Migration file not found: $MIGRATION_FILE"

log "Checking schema state..."
SHOPS_TABLE_EXISTS="$(psql "$APP_DB_CONN" -tAc "SELECT to_regclass('public.shops') IS NOT NULL" | tr -d '[:space:]')"
if [[ "$SHOPS_TABLE_EXISTS" != "t" ]]; then
  log "Applying initial migration..."
  psql "$APP_DB_CONN" -f "$MIGRATION_FILE" >/dev/null
  log "Applying categories migration..."
  psql "$APP_DB_CONN" -f "$ROOT_DIR/migrations/000003_add_categories.up.sql" >/dev/null
  log "Migrations applied"
else
  log "Schema already initialized, skipping migration"
fi

SUPERADMIN_EXISTS="$(psql "$APP_DB_CONN" -tAc "SELECT EXISTS (SELECT 1 FROM users WHERE role = 'SuperAdmin')" | tr -d '[:space:]')"
if [[ "$SUPERADMIN_EXISTS" != "t" ]]; then
  log "No SuperAdmin detected, preparing initial seed..."

  : "${FIRST_SHOP_ID:=$(generate_uuid)}"
  : "${FIRST_SHOP_NAME:=Shop Demo}"
  : "${FIRST_SHOP_WHATSAPP:=+212600000000}"
  : "${FIRST_SUPERADMIN_NAME:=Owner Demo}"
  : "${FIRST_SUPERADMIN_EMAIL:=owner@shopdemo.com}"

  if [[ -z "${FIRST_SUPERADMIN_PASSWORD:-}" ]]; then
    if [[ -t 0 ]]; then
      read -r -s -p "Enter FIRST_SUPERADMIN_PASSWORD: " FIRST_SUPERADMIN_PASSWORD
      echo
    else
      fail "FIRST_SUPERADMIN_PASSWORD is required in non-interactive mode"
    fi
  fi

  [[ -n "$FIRST_SUPERADMIN_PASSWORD" ]] || fail "FIRST_SUPERADMIN_PASSWORD cannot be empty"
  [[ "$FIRST_SHOP_WHATSAPP" =~ ^\+[1-9][0-9]{7,14}$ ]] || fail "FIRST_SHOP_WHATSAPP must match +[country][number], e.g. +212600000000"
  [[ "$BCRYPT_COST" =~ ^[0-9]+$ ]] || fail "BCRYPT_COST must be numeric"

  psql "$APP_DB_CONN" \
    -v shop_id="$FIRST_SHOP_ID" \
    -v shop_name="$FIRST_SHOP_NAME" \
    -v shop_whatsapp="$FIRST_SHOP_WHATSAPP" \
    -v admin_name="$FIRST_SUPERADMIN_NAME" \
    -v admin_email="$FIRST_SUPERADMIN_EMAIL" \
    -v admin_password="$FIRST_SUPERADMIN_PASSWORD" \
    -v bcrypt_cost="$BCRYPT_COST" <<'SQL' >/dev/null
INSERT INTO shops (id, name, active, whatsapp_number)
VALUES (:'shop_id'::uuid, :'shop_name', true, :'shop_whatsapp')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (name, email, password, role, shop_id)
VALUES (
  :'admin_name',
  LOWER(:'admin_email'),
  crypt(:'admin_password', gen_salt('bf', :'bcrypt_cost'::int)),
  'SuperAdmin',
  :'shop_id'::uuid
)
ON CONFLICT DO NOTHING;
SQL

  log "Initial seed created"
  log "  shopID:   ${FIRST_SHOP_ID}"
  log "  email:    ${FIRST_SUPERADMIN_EMAIL}"
  log "  password: (the one you entered)"
else
  log "SuperAdmin already exists, skipping seed"
fi

if [[ "$NO_RUN" == "1" ]]; then
  log "Preparation complete. Start API with: go run ./cmd/api"
  exit 0
fi

log "Starting API..."
exec go run ./cmd/api
