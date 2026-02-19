# scripts/docker-test.ps1 - Full API smoke test (Windows PowerShell)
# Usage: .\scripts\docker-test.ps1

$ErrorActionPreference = "Stop"

$BASE_URL = "http://localhost:8080"
$SHOP_ID = "11111111-1111-1111-1111-111111111111"
$EMAIL = "owner@shopdemo.com"
$PASSWORD = "ChangeMe123!"

$pass = 0
$fail = 0

function Green($msg) { Write-Host "  [OK] $msg" -ForegroundColor Green }
function Red($msg) { Write-Host "  [KO] $msg" -ForegroundColor Red }
function Blue($msg) { Write-Host $msg -ForegroundColor Cyan }

function Check($label, $expected, $actual) {
    if ($actual -eq $expected) {
        Green $label
        $script:pass++
    } else {
        Red "$label (expected: $expected, actual: $actual)"
        $script:fail++
    }
}

function CheckNotEmpty($label, $value) {
    if ($value -and $value -ne "null") {
        Green "$label ($value)"
        $script:pass++
    } else {
        Red "$label (empty or null value)"
        $script:fail++
    }
}

Blue "============================================================"
Blue "  Arem-Shop - Full API Docker Test (Windows)"
Blue "  $BASE_URL"
Blue "============================================================"
Write-Host ""

# 1. Healthcheck
Blue "1) Healthcheck"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/health" -UseBasicParsing
    Check "GET /health -> 200" 200 $r.StatusCode
} catch {
    Red "GET /health -> Error (API not reachable?)"
    $fail++
}
Write-Host ""

# 2. Login SuperAdmin
Blue "2) Login SuperAdmin"
$body = @{
    email = $EMAIL
    password = $PASSWORD
    shopID = $SHOP_ID
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/auth/login" -Method POST `
        -ContentType "application/json" -Body $body -UseBasicParsing
    $token = ($r.Content | ConvertFrom-Json).token
    CheckNotEmpty "JWT token retrieved" $token
} catch {
    Red "Login failed: $($_.Exception.Message)"
    $fail++
    Write-Host "Cannot continue without token." -ForegroundColor Red
    exit 1
}

if (-not $token) {
    Write-Host "Cannot continue without token." -ForegroundColor Red
    exit 1
}

$headers = @{ Authorization = "Bearer $token" }
Write-Host ""

# 3. Create a product
Blue "3) Create a product (SuperAdmin)"
$prodBody = @{
    name = "iPhone 14 Test"
    description = "128GB - Smoke Test"
    category = "Smartphones"
    purchasePrice = 4500.00
    sellingPrice = 5999.99
    stock = 10
    imageURL = "https://example.com/iphone14.jpg"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Method POST `
        -ContentType "application/json" -Body $prodBody -Headers $headers -UseBasicParsing
    $productId = ($r.Content | ConvertFrom-Json).id
    CheckNotEmpty "Product created with ID" $productId
} catch {
    Red "Product creation failed: $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# 4. List private products
Blue "4) List private products"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Headers $headers -UseBasicParsing
    Check "GET /products -> 200" 200 $r.StatusCode
} catch {
    Red "GET /products failed"
    $fail++
}
Write-Host ""

# 5. Public catalog
Blue "5) Public catalog"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/public/$SHOP_ID/products" -UseBasicParsing
    Check "GET /public/:shopID/products -> 200" 200 $r.StatusCode

    if ($r.Content -match "purchasePrice") {
        Red "purchasePrice is visible in public catalog"
        $fail++
    } else {
        Green "purchasePrice is hidden in public catalog"
        $pass++
    }
} catch {
    Red "Public catalog failed: $($_.Exception.Message)"
    $fail += 2
}
Write-Host ""

# 6. Create + delete a product (no transactions)
Blue "6) Create + delete a product (no transactions)"
$tempBody = @{
    name = "Temp Product"
    description = "Delete test"
    category = "Tests"
    purchasePrice = 10.00
    sellingPrice = 20.00
    stock = 1
    imageURL = "https://example.com/temp.jpg"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/products" -Method POST `
        -ContentType "application/json" -Body $tempBody -Headers $headers -UseBasicParsing
    $tempId = ($r.Content | ConvertFrom-Json).id

    $r = Invoke-WebRequest -Uri "$BASE_URL/products/$tempId" -Method DELETE `
        -Headers $headers -UseBasicParsing
    Check "DELETE /products/:id -> 204" 204 $r.StatusCode
} catch {
    Red "Delete failed: $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# 7. Sale transaction
Blue "7) Sale transaction (stock: 10 -> 8)"
$saleBody = @{
    type = "Sale"
    productID = $productId
    quantity = 2
    amount = 11999.98
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/transactions" -Method POST `
        -ContentType "application/json" -Body $saleBody -Headers $headers -UseBasicParsing
    Check "POST /transactions (Sale) -> 201" 201 $r.StatusCode
} catch {
    Red "Sale failed: $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# 8. Expense transaction
Blue "8) Expense transaction"
$expBody = @{
    type = "Expense"
    amount = 500.00
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/transactions" -Method POST `
        -ContentType "application/json" -Body $expBody -Headers $headers -UseBasicParsing
    Check "POST /transactions (Expense) -> 201" 201 $r.StatusCode
} catch {
    Red "Expense failed: $($_.Exception.Message)"
    $fail++
}
Write-Host ""

# 9. SuperAdmin dashboard
Blue "9) SuperAdmin dashboard"
try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/reports/dashboard" -Headers $headers -UseBasicParsing
    Check "GET /reports/dashboard -> 200" 200 $r.StatusCode
} catch {
    Red "Dashboard failed"
    $fail++
}
Write-Host ""

# 10. Security: cross-shop
Blue "10) Security: cross-shop attempt"
$crossBody = @{
    name = "Cross Shop"
    email = "cross@test.com"
    password = "StrongPass123"
    role = "Admin"
    shopID = "22222222-2222-2222-2222-222222222222"
} | ConvertTo-Json

try {
    $r = Invoke-WebRequest -Uri "$BASE_URL/auth/register" -Method POST `
        -ContentType "application/json" -Body $crossBody -Headers $headers -UseBasicParsing
    Red "Cross-shop accepted (should be rejected)"
    $fail++
} catch {
    $code = $_.Exception.Response.StatusCode.value__
    Check "POST /auth/register cross-shop -> 403" 403 $code
}
Write-Host ""

# Result
Blue "============================================================"
$total = $pass + $fail
if ($fail -eq 0) {
    Write-Host "  ALL TESTS PASSED: $pass/$total" -ForegroundColor Green
} else {
    Write-Host "  RESULT: $pass OK / $fail FAILURES" -ForegroundColor Red
}
Blue "============================================================"

exit $fail
