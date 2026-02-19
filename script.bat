@echo off
setlocal EnableDelayedExpansion

rem ============================================================
rem Arem Shop API - Test all endpoints
rem ============================================================

set "BASE_URL=https://special-barnacle-g4xjw6qv76wvf9pj4-8080.app.github.dev"
set "SHOP_ID=11111111-1111-1111-1111-111111111111"
set "SUPERADMIN_EMAIL=owner@shopdemo.com"
set "SUPERADMIN_PASSWORD=ChangeMe123!"
set "NEW_ADMIN_EMAIL=admin_%RANDOM%%RANDOM%@shopdemo.com"
set "NEW_ADMIN_PASSWORD=AdminPass123!"

set "WORKDIR=%TEMP%\arem-shop-api-test"
if not exist "%WORKDIR%" mkdir "%WORKDIR%"

set "TOKEN="
set "PRODUCT_ID="
set "TOKEN_FILE=%WORKDIR%\token.txt"
set "PRODUCT_ID_FILE=%WORKDIR%\product_id.txt"

echo.
echo [1/11] GET /health
curl -s -X GET "%BASE_URL%/health" -H "Accept: application/json" -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [2/11] POST /auth/login
(
echo {
echo   "email":"%SUPERADMIN_EMAIL%",
echo   "password":"%SUPERADMIN_PASSWORD%",
echo   "shopID":"%SHOP_ID%"
echo }
)> "%WORKDIR%\login.json"

curl -s -X POST "%BASE_URL%/auth/login" ^
  -H "Content-Type: application/json" ^
  --data "@%WORKDIR%\login.json" ^
  -o "%WORKDIR%\login_resp.json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

if exist "%TOKEN_FILE%" del "%TOKEN_FILE%"
powershell -NoProfile -Command ^
  "$j = Get-Content -Raw '%WORKDIR%\login_resp.json' | ConvertFrom-Json; if($j.token){[System.IO.File]::WriteAllText('%TOKEN_FILE%', $j.token)}"

if exist "%TOKEN_FILE%" set /p TOKEN=<"%TOKEN_FILE%"

if "%TOKEN%"=="" (
  echo.
  echo ERROR: token not found in login response. Stopping.
  type "%WORKDIR%\login_resp.json"
  exit /b 1
)

echo Token loaded.

echo.
echo [3/11] POST /auth/register
(
echo {
echo   "name":"Admin Test",
echo   "email":"%NEW_ADMIN_EMAIL%",
echo   "password":"%NEW_ADMIN_PASSWORD%",
echo   "role":"Admin",
echo   "shopID":"%SHOP_ID%"
echo }
)> "%WORKDIR%\register.json"

curl -s -X POST "%BASE_URL%/auth/register" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\register.json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [4/11] POST /products
(
echo {
echo   "name":"Laptop Pro Test",
echo   "description":"Produit cree par script.bat",
echo   "category":"Laptops",
echo   "purchasePrice":"800.00",
echo   "sellingPrice":"999.99",
echo   "stock":15,
echo   "imageURL":"https://example.com/laptop.jpg"
echo }
)> "%WORKDIR%\product_create.json"

curl -s -X POST "%BASE_URL%/products" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\product_create.json" ^
  -o "%WORKDIR%\product_create_resp.json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

if exist "%PRODUCT_ID_FILE%" del "%PRODUCT_ID_FILE%"
powershell -NoProfile -Command ^
  "$j = Get-Content -Raw '%WORKDIR%\product_create_resp.json' | ConvertFrom-Json; $id=''; if($j.data -and $j.data.id){$id=$j.data.id} elseif($j.id){$id=$j.id}; if($id){[System.IO.File]::WriteAllText('%PRODUCT_ID_FILE%', $id)}"

if exist "%PRODUCT_ID_FILE%" set /p PRODUCT_ID=<"%PRODUCT_ID_FILE%"

if "%PRODUCT_ID%"=="" (
  echo.
  echo ERROR: product id not found in product create response. Stopping.
  type "%WORKDIR%\product_create_resp.json"
  exit /b 1
)

echo Product ID: %PRODUCT_ID%

echo.
echo [5/11] GET /products
curl -s -X GET "%BASE_URL%/products" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [6/11] GET /products/%PRODUCT_ID%
curl -s -X GET "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [7/11] PUT /products/%PRODUCT_ID%
(
echo {
echo   "name":"Laptop Pro Test Updated",
echo   "description":"Produit modifie par script.bat",
echo   "category":"Laptops",
echo   "purchasePrice":"790.00",
echo   "sellingPrice":"1099.99",
echo   "stock":12,
echo   "imageURL":"https://example.com/laptop-updated.jpg"
echo }
)> "%WORKDIR%\product_update.json"

curl -s -X PUT "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\product_update.json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [8/11] GET /public/%SHOP_ID%/products
curl -s -X GET "%BASE_URL%/public/%SHOP_ID%/products" ^
  -H "Accept: application/json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [9/11] POST /transactions (Sale)
(
echo {
echo   "type":"Sale",
echo   "productID":"%PRODUCT_ID%",
echo   "quantity":1,
echo   "amount":"1099.99"
echo }
)> "%WORKDIR%\transaction.json"

curl -s -X POST "%BASE_URL%/transactions" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\transaction.json" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [10/11] GET /reports/dashboard
curl -s -X GET "%BASE_URL%/reports/dashboard" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo [11/11] DELETE /products/%PRODUCT_ID%
curl -s -X DELETE "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -w "\nHTTP_STATUS:%%{http_code}\n"

echo.
echo Tests finished.
echo Temp files: %WORKDIR%
exit /b 0

