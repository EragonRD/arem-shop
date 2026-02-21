@echo off
setlocal EnableExtensions DisableDelayedExpansion

rem ============================================================
rem Arem Shop API - Endpoint test script (colorized logs)
rem ============================================================

set "BASE_URL=https://laughing-parakeet-jjrvw95pq959fpggw-8080.app.github.dev"
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
set "HTTP_CODE="
set /a PASS_COUNT=0
set /a FAIL_COUNT=0
set "ABORTED=0"

for /f %%E in ('echo prompt $E^| cmd') do set "ESC=%%E"
if defined ESC (
  set "CLR_RESET=%ESC%[0m"
  set "CLR_BOLD=%ESC%[1m"
  set "CLR_CYAN=%ESC%[36m"
  set "CLR_GREEN=%ESC%[32m"
  set "CLR_YELLOW=%ESC%[33m"
  set "CLR_RED=%ESC%[31m"
) else (
  set "CLR_RESET="
  set "CLR_BOLD="
  set "CLR_CYAN="
  set "CLR_GREEN="
  set "CLR_YELLOW="
  set "CLR_RED="
)

call :banner "Arem Shop API Smoke Test"
call :info "Base URL : %BASE_URL%"
call :info "Shop ID  : %SHOP_ID%"

call :step "1/11" "GET /health"
curl -sS -X GET "%BASE_URL%/health" ^
  -H "Accept: application/json" ^
  -o "%WORKDIR%\health_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\health_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\health_status.txt"
type "%WORKDIR%\health_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "GET /health"

call :step "2/11" "POST /auth/login"
(
echo {
echo   "email":"%SUPERADMIN_EMAIL%",
echo   "password":"%SUPERADMIN_PASSWORD%",
echo   "shopID":"%SHOP_ID%"
echo }
)> "%WORKDIR%\login.json"

curl -sS -X POST "%BASE_URL%/auth/login" ^
  -H "Content-Type: application/json" ^
  --data "@%WORKDIR%\login.json" ^
  -o "%WORKDIR%\login_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\login_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\login_status.txt"
type "%WORKDIR%\login_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "POST /auth/login"

if not "%HTTP_CODE%"=="200" (
  call :fatal "Login failed. Check SUPERADMIN_EMAIL/SUPERADMIN_PASSWORD/SHOP_ID."
  goto :script_end
)

if exist "%TOKEN_FILE%" del "%TOKEN_FILE%"
powershell -NoProfile -Command ^
  "$j = Get-Content -Raw '%WORKDIR%\login_resp.json' | ConvertFrom-Json; if($j.token){[System.IO.File]::WriteAllText('%TOKEN_FILE%', $j.token)}"

if exist "%TOKEN_FILE%" set /p TOKEN=<"%TOKEN_FILE%"
if "%TOKEN%"=="" (
  call :fatal "Token not found in login response."
  goto :script_end
)
call :ok "Token loaded."

call :step "3/11" "POST /auth/register"
(
echo {
echo   "name":"Admin Test",
echo   "email":"%NEW_ADMIN_EMAIL%",
echo   "password":"%NEW_ADMIN_PASSWORD%",
echo   "role":"Admin",
echo   "shopID":"%SHOP_ID%"
echo }
)> "%WORKDIR%\register.json"

curl -sS -X POST "%BASE_URL%/auth/register" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\register.json" ^
  -o "%WORKDIR%\register_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\register_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\register_status.txt"
type "%WORKDIR%\register_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "201" "POST /auth/register"

call :step "4/11" "POST /products"
(
echo {
echo   "name":"Laptop Pro Test",
echo   "description":"Produit cree par script.bat",
echo   "category":"Laptops",
echo   "purchasePrice":800.00,
echo   "sellingPrice":999.99,
echo   "stock":15,
echo   "imageURL":"https://example.com/laptop.jpg"
echo }
)> "%WORKDIR%\product_create.json"

curl -sS -X POST "%BASE_URL%/products" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\product_create.json" ^
  -o "%WORKDIR%\product_create_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\product_create_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\product_create_status.txt"
type "%WORKDIR%\product_create_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "201" "POST /products"

if not "%HTTP_CODE%"=="201" (
  call :fatal "Product creation failed. Cannot continue scenario."
  goto :script_end
)

if exist "%PRODUCT_ID_FILE%" del "%PRODUCT_ID_FILE%"
powershell -NoProfile -Command ^
  "$j = Get-Content -Raw '%WORKDIR%\product_create_resp.json' | ConvertFrom-Json; $id=''; if($j.data -and $j.data.id){$id=$j.data.id} elseif($j.id){$id=$j.id}; if($id){[System.IO.File]::WriteAllText('%PRODUCT_ID_FILE%', $id)}"

if exist "%PRODUCT_ID_FILE%" set /p PRODUCT_ID=<"%PRODUCT_ID_FILE%"
if "%PRODUCT_ID%"=="" (
  call :fatal "Product ID not found in create response."
  goto :script_end
)
call :ok "Product ID: %PRODUCT_ID%"

call :step "5/11" "GET /products"
curl -sS -X GET "%BASE_URL%/products" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -o "%WORKDIR%\products_list_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\products_list_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\products_list_status.txt"
type "%WORKDIR%\products_list_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "GET /products"

call :step "6/11" "GET /products/%PRODUCT_ID%"
curl -sS -X GET "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -o "%WORKDIR%\product_get_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\product_get_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\product_get_status.txt"
type "%WORKDIR%\product_get_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "GET /products/:id"

call :step "7/11" "PUT /products/%PRODUCT_ID%"
(
echo {
echo   "name":"Laptop Pro Test Updated",
echo   "description":"Produit modifie par script.bat",
echo   "category":"Laptops",
echo   "purchasePrice":790.00,
echo   "sellingPrice":1099.99,
echo   "stock":12,
echo   "imageURL":"https://example.com/laptop-updated.jpg"
echo }
)> "%WORKDIR%\product_update.json"

curl -sS -X PUT "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\product_update.json" ^
  -o "%WORKDIR%\product_update_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\product_update_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\product_update_status.txt"
type "%WORKDIR%\product_update_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "PUT /products/:id"

call :step "8/11" "GET /public/%SHOP_ID%/products"
curl -sS -X GET "%BASE_URL%/public/%SHOP_ID%/products" ^
  -H "Accept: application/json" ^
  -o "%WORKDIR%\public_products_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\public_products_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\public_products_status.txt"
type "%WORKDIR%\public_products_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "GET /public/:shopID/products"

call :step "9/11" "POST /transactions (Sale)"
(
echo {
echo   "type":"Sale",
echo   "productID":"%PRODUCT_ID%",
echo   "quantity":1,
echo   "amount":1099.99
echo }
)> "%WORKDIR%\transaction.json"

curl -sS -X POST "%BASE_URL%/transactions" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  --data "@%WORKDIR%\transaction.json" ^
  -o "%WORKDIR%\transaction_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\transaction_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\transaction_status.txt"
type "%WORKDIR%\transaction_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "201" "POST /transactions"

call :step "10/11" "GET /reports/dashboard"
curl -sS -X GET "%BASE_URL%/reports/dashboard" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -o "%WORKDIR%\dashboard_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\dashboard_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\dashboard_status.txt"
type "%WORKDIR%\dashboard_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200" "GET /reports/dashboard"

call :step "11/11" "DELETE /products/%PRODUCT_ID%"
curl -sS -X DELETE "%BASE_URL%/products/%PRODUCT_ID%" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -o "%WORKDIR%\product_delete_resp.json" ^
  -w "%%{http_code}" > "%WORKDIR%\product_delete_status.txt"
set /p HTTP_CODE=<"%WORKDIR%\product_delete_status.txt"
type "%WORKDIR%\product_delete_resp.json"
echo.
call :print_http_status "%HTTP_CODE%"
call :check_status "%HTTP_CODE%" "200,409" "DELETE /products/:id"

:script_end
echo.
echo %CLR_BOLD%%CLR_CYAN%============================================================%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%Summary%CLR_RESET%
echo %CLR_GREEN%Passed: %PASS_COUNT%%CLR_RESET%
echo %CLR_RED%Failed: %FAIL_COUNT%%CLR_RESET%
if "%ABORTED%"=="1" echo %CLR_RED%Execution aborted early due to blocking failure.%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%Temp files: %WORKDIR%%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%============================================================%CLR_RESET%
echo.

if %FAIL_COUNT% gtr 0 (
  exit /b 1
)
exit /b 0

:banner
echo.
echo %CLR_BOLD%%CLR_CYAN%============================================================%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%%~1%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%============================================================%CLR_RESET%
echo.
exit /b 0

:step
echo.
echo %CLR_BOLD%%CLR_CYAN%[%~1] %~2%CLR_RESET%
exit /b 0

:info
echo %CLR_CYAN%[INFO] %~1%CLR_RESET%
exit /b 0

:ok
echo %CLR_GREEN%[OK] %~1%CLR_RESET%
exit /b 0

:warn
echo %CLR_YELLOW%[WARN] %~1%CLR_RESET%
exit /b 0

:print_http_status
set "STATUS_COLOR=%CLR_YELLOW%"
if "%~1"=="" (
  set "STATUS_VALUE=000"
) else (
  set "STATUS_VALUE=%~1"
)
if "%STATUS_VALUE:~0,1%"=="2" set "STATUS_COLOR=%CLR_GREEN%"
if "%STATUS_VALUE:~0,1%"=="4" set "STATUS_COLOR=%CLR_YELLOW%"
if "%STATUS_VALUE:~0,1%"=="5" set "STATUS_COLOR=%CLR_RED%"
echo %STATUS_COLOR%HTTP_STATUS:%STATUS_VALUE%%CLR_RESET%
exit /b 0

:check_status
set "ACTUAL=%~1"
set "EXPECTED=%~2"
set "LABEL=%~3"
echo ,%EXPECTED%, | findstr /i /c:",%ACTUAL%," >nul
if errorlevel 1 (
  set /a FAIL_COUNT+=1
  call :warn "%LABEL% -> got %ACTUAL% (expected %EXPECTED%)"
  exit /b 1
)
set /a PASS_COUNT+=1
call :ok "%LABEL%"
exit /b 0

:fatal
set /a FAIL_COUNT+=1
set "ABORTED=1"
echo.
echo %CLR_RED%[ERROR] %~1%CLR_RESET%
echo %CLR_BOLD%%CLR_CYAN%Temp files: %WORKDIR%%CLR_RESET%
exit /b 0

