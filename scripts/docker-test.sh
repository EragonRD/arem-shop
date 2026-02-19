#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# scripts/docker-test.sh — Test complet de l'API après docker-compose up
# ──────────────────────────────────────────────────────────────
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# ── Config ────────────────────────────────────────────────────
BASE_URL="${BASE_URL:-http://localhost:8080}"
SHOP_ID="11111111-1111-1111-1111-111111111111"
SUPERADMIN_EMAIL="owner@shopdemo.com"
SUPERADMIN_PASSWORD="ChangeMe123!"

PASS=0
FAIL=0

green() { printf '\033[1;32m%s\033[0m\n' "$*"; }
red()   { printf '\033[1;31m%s\033[0m\n' "$*"; }
blue()  { printf '\033[1;34m%s\033[0m\n' "$*"; }

check() {
  local label="$1" expected="$2" actual="$3"
  if [[ "$actual" == "$expected" ]]; then
    green "  ✅ $label"
    PASS=$((PASS + 1))
  else
    red "  ❌ $label (attendu: $expected, obtenu: $actual)"
    FAIL=$((FAIL + 1))
  fi
}

check_not_empty() {
  local label="$1" actual="$2"
  if [[ -n "$actual" && "$actual" != "null" ]]; then
    green "  ✅ $label ($actual)"
    PASS=$((PASS + 1))
  else
    red "  ❌ $label (valeur vide ou null)"
    FAIL=$((FAIL + 1))
  fi
}

blue "══════════════════════════════════════════════════════════"
blue "  Arem-Shop — Test complet API Docker"
blue "  $BASE_URL"
blue "══════════════════════════════════════════════════════════"
echo

# ── 1. Healthcheck ────────────────────────────────────────────
blue "1) Healthcheck"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
check "GET /health → 200" "200" "$HTTP_CODE"
echo

# ── 2. Login SuperAdmin ──────────────────────────────────────
blue "2) Login SuperAdmin"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$SUPERADMIN_EMAIL\",
    \"password\": \"$SUPERADMIN_PASSWORD\",
    \"shopID\": \"$SHOP_ID\"
  }")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
check_not_empty "Token JWT obtenu" "$TOKEN"

if [[ -z "${TOKEN:-}" ]]; then
  red "Impossible de continuer sans token. Réponse login :"
  echo "$LOGIN_RESPONSE" | jq . 2>/dev/null || echo "$LOGIN_RESPONSE"
  exit 1
fi
echo

# ── 3. Créer un produit (SuperAdmin) ─────────────────────────
blue "3) Créer un produit (SuperAdmin)"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 14 Test",
    "description": "128GB - Smoke Test",
    "category": "Smartphones",
    "purchasePrice": 4500.00,
    "sellingPrice": 5999.99,
    "stock": 10,
    "imageURL": "https://example.com/iphone14.jpg"
  }')

PRODUCT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id // empty')
check_not_empty "Produit créé avec ID" "$PRODUCT_ID"
echo

# ── 4. Lister les produits privés ─────────────────────────────
blue "4) Lister les produits privés"
LIST_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" "$BASE_URL/products")
check "GET /products → 200" "200" "$LIST_CODE"
echo

# ── 5. Catalogue public ──────────────────────────────────────
blue "5) Catalogue public"
PUBLIC_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/public/$SHOP_ID/products")
check "GET /public/:shopID/products → 200" "200" "$PUBLIC_CODE"

# Vérifier que purchasePrice n'est pas exposé (cherche dans le JSON brut)
PUBLIC_BODY=$(curl -s "$BASE_URL/public/$SHOP_ID/products")
if echo "$PUBLIC_BODY" | grep -qi "purchasePrice"; then
  red "  ❌ purchasePrice visible dans le catalogue public"
  FAIL=$((FAIL + 1))
else
  green "  ✅ purchasePrice masqué dans le catalogue public"
  PASS=$((PASS + 1))
fi
echo

# ── 6. Suppression produit (avant transactions) ──────────────
# On crée un 2e produit temporaire pour tester le DELETE
blue "6) Créer + supprimer un produit (sans transactions)"
TEMP_RESPONSE=$(curl -s -X POST "$BASE_URL/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produit Temp",
    "description": "Pour test delete",
    "category": "Tests",
    "purchasePrice": 10.00,
    "sellingPrice": 20.00,
    "stock": 1,
    "imageURL": "https://example.com/temp.jpg"
  }')
TEMP_ID=$(echo "$TEMP_RESPONSE" | jq -r '.id // empty')
if [[ -n "$TEMP_ID" && "$TEMP_ID" != "null" ]]; then
  DEL_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE \
    "$BASE_URL/products/$TEMP_ID" \
    -H "Authorization: Bearer $TOKEN")
  check "DELETE /products/:id → 204" "204" "$DEL_CODE"
else
  red "  ❌ Impossible de créer le produit temporaire"
  FAIL=$((FAIL + 1))
fi
echo

# ── 7. Transaction Sale (décrément stock) ────────────────────
blue "7) Transaction Sale (stock: 10 → 8)"
SALE_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/transactions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"Sale\",
    \"productID\": \"$PRODUCT_ID\",
    \"quantity\": 2,
    \"amount\": 11999.98
  }")
check "POST /transactions (Sale) → 201" "201" "$SALE_CODE"
echo

# ── 8. Transaction Expense ────────────────────────────────────
blue "8) Transaction Expense"
EXPENSE_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/transactions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Expense",
    "amount": 500.00
  }')
check "POST /transactions (Expense) → 201" "201" "$EXPENSE_CODE"
echo

# ── 9. Dashboard SuperAdmin ──────────────────────────────────
blue "9) Dashboard SuperAdmin"
DASH_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" "$BASE_URL/reports/dashboard")
check "GET /reports/dashboard → 200" "200" "$DASH_CODE"
echo

# ── 10. Sécurité : cross-shop interdit ────────────────────────
blue "10) Sécurité : tentative cross-shop"
CROSS_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/auth/register" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cross Shop",
    "email": "cross@test.com",
    "password": "StrongPass123",
    "role": "Admin",
    "shopID": "22222222-2222-2222-2222-222222222222"
  }')
check "POST /auth/register cross-shop → 403" "403" "$CROSS_CODE"
echo

# ── Résultat ──────────────────────────────────────────────────
blue "══════════════════════════════════════════════════════════"
if [[ "$FAIL" -eq 0 ]]; then
  green "  🎉 TOUS LES TESTS PASSENT : $PASS/$((PASS + FAIL))"
else
  red   "  ⚠️  RÉSULTAT : $PASS OK / $FAIL ÉCHECS"
fi
blue "══════════════════════════════════════════════════════════"

exit "$FAIL"
