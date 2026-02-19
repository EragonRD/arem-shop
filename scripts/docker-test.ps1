# ──────────────────────────────────────────────────────────────
# scripts/docker-test.ps1 — Test complet de l'API (Windows PowerShell)
# Usage : .\scripts\docker-test.ps1
# ──────────────────────────────────────────────────────────────

$ErrorActionPreference = "Stop"

$BASE_URL      = "http://localhost:8080"
$SHOP_ID       = "11111111-1111-1111-1111-111111111111"
$EMAIL         = "owner@shopdemo.com"
$PASSWORD      = "ChangeMe123!"

$pass = 0
$fail = 0

function Green  ($msg) { Write-Host "  ✅ $msg" -ForegroundColor Green }
function Red    ($msg) { Write-Host "  ❌ $msg" -ForegroundColor Red }
function Blue   ($msg) { Write-Host $msg -ForegroundColor Cyan }

function Check ($label, $expected, $actual) {
    if ($actual -eq $expected) {
        Green $label
        $script:pass++
    } else {
        Red "$label (attendu: $expected, obtenu: $actual)"
        $script:fail++
    }
}

function CheckNotEmpty ($label, $value) {
    if ($value -and $value -ne "null") {
        Green "$label ($value)"
        $script:pass++
    } else {
        Red "$label (valeur vide ou null)"
        $script:fail++
    }
}

Blue "══════════════════════════════════════════════════════════"
Blue "  Arem-Shop — Test complet API Docker (Windows)"
Blue "  $BASE_URL"
Blue "══════════════════════════════════════════════════════════"
Write-Host ""

# ── 1. Healthcheck ────────────────────────────────────────────
Blue "1) Healthcheck"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/health" -UseBasicParsing
    Check "GET /health → 200" 200 $r.StatusCode
} catch {
    Red "GET /health → Erreur (API non accessible ?)"
    $fail++
}
Write-Host ""

# ── 2. Login SuperAdmin ──────────────────────────────────────
Blue "2) Login SuperAdmin"
$body = @{
    email    = $EMAIL
    password = $PASSWORD
    shopID   = $SHOP_ID
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/auth/login" -Method POST `
        -ContentType "application/json" -Body $body -UseBasicParsing
    $token = ($r.Content | ConvertFrom-Json).token
    CheckNotEmpty "Token JWT obtenu" $token
} catch {
    Red "Login echoue : $($_.Exception.Message)"
    $fail++
    Write-Host "Impossible de continuer sans token." -ForegroundColor Red
    exit 1
}

if (-not $token) {
    Write-Host "Impossible de continuer sans token." -ForegroundColor Red
    exit 1
}

$headers = @{ Authorization = "Bearer $token" }
Write-Host ""

# ── 3. Creer un produit ──────────────────────────────────────
Blue "3) Creer un produit (SuperAdmin)"
$prodBody = @{
    name          = "iPhone 14 Test"
    description   = "128GB - Smoke Test"
    category      = "Smartphones"
    purchasePrice = 4500.00
    sellingPrice  = 5999.99
    stock         = 10
    imageURL      = "https://example.com/iphone14.jpg"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Method POST `
        -ContentType "application/json" -Body $prodBody -Headers $headers -UseBasicParsing
    $productId = ($r.Content | ConvertFrom-Json).id
    CheckNotEmpty "Produit cree avec ID" $productId
} catch {
    Red "Creation produit echouee : $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# ── 4. Lister les produits prives ─────────────────────────────
Blue "4) Lister les produits prives"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Headers $headers -UseBasicParsing
    Check "GET /products → 200" 200 $r.StatusCode
} catch {
    Red "GET /products echoue"
    $fail++
}
Write-Host ""

# ── 5. Catalogue public ──────────────────────────────────────
Blue "5) Catalogue public"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/public/$SHOP_ID/products" -UseBasicParsing
    Check "GET /public/:shopID/products → 200" 200 $r.StatusCode

    if ($r.Content -match "purchasePrice") {
        Red "purchasePrice visible dans le catalogue public"
        $fail++
    } else {
        Green "purchasePrice masque dans le catalogue public"
        $pass++
    }
} catch {
    Red "Catalogue public echoue : $($_.Exception.Message)"
    $fail += 2
}
Write-Host ""

# ── 6. Supprimer un produit (sans transactions) ──────────────
Blue "6) Creer + supprimer un produit (sans transactions)"
$tempBody = @{
    name          = "Produit Temp"
    description   = "Pour test delete"
    category      = "Tests"
    purchasePrice = 10.00
    sellingPrice  = 20.00
    stock         = 1
    imageURL      = "https://example.com/temp.jpg"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Method POST `
        -ContentType "application/json" -Body $tempBody -Headers $headers -UseBasicParsing
    $tempId = ($r.Content | ConvertFrom-Json).id

    $r = Invoke-WebRequest -Uri "$BASE_URL/products/$tempId" -Method DELETE `
        -Headers $headers -UseBasicParsing
    Check "DELETE /products/:id → 204" 204 $r.StatusCode
} catch {
    Red "Delete echoue : $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# ── 7. Transaction Sale ──────────────────────────────────────
Blue "7) Transaction Sale (stock: 10 → 8)"
$saleBody = @{
    type      = "Sale"
    productID = $productId
    quantity  = 2
    amount    = 11999.98
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/transactions" -Method POST `
        -ContentType "application/json" -Body $saleBody -Headers $headers -UseBasicParsing
    Check "POST /transactions (Sale) → 201" 201 $r.StatusCode
} catch {
    Red "Sale echouee : $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# ── 8. Transaction Expense ────────────────────────────────────
Blue "8) Transaction Expense"
$expBody = @{
    type   = "Expense"
    amount = 500.00
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/transactions" -Method POST `
        -ContentType "application/json" -Body $expBody -Headers $headers -UseBasicParsing
    Check "POST /transactions (Expense) → 201" 201 $r.StatusCode
} catch {
    Red "Expense echouee : $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# ── 9. Dashboard SuperAdmin ──────────────────────────────────
Blue "9) Dashboard SuperAdmin"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/reports/dashboard" -Headers $headers -UseBasicParsing
    Check "GET /reports/dashboard → 200" 200 $r.StatusCode
} catch {
    Red "Dashboard echoue"
    $fail++
}
Write-Host ""

# ── 10. Securite : cross-shop ────────────────────────────────
Blue "10) Securite : tentative cross-shop"
$crossBody = @{
    name     = "Cross Shop"
    email    = "cross@test.com"
    password = "StrongPass123"
    role     = "Admin"
    shopID   = "22222222-2222-2222-2222-222222222222"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/auth/register" -Method POST `
        -ContentType "application/json" -Body $crossBody -Headers $headers -UseBasicParsing
    Red "Cross-shop accepte (devrait etre refuse)"
    $fail++
} catch {
    $code = $_.Exception.Response.StatusCode.value__
    Check "POST /auth/register cross-shop → 403" 403 $code
}
Write-Host ""

# ── Resultat ──────────────────────────────────────────────────
Blue "══════════════════════════════════════════════════════════"
$total = $pass + $fail
if ($fail -eq 0) {
    Write-Host "  🎉 TOUS LES TESTS PASSENT : $pass/$total" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  RESULTAT : $pass OK / $fail ECHECS" -ForegroundColor Red
}
Blue "══════════════════════════════════════════════════════════"

exit $fail
