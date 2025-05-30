@echo off
REM Настройки SMTP для Gmail
set SMTP_HOST=smtp.gmail.com
set SMTP_PORT=587
REM Gmail аккаунт
set SMTP_USER=kashkenov2006@gmail.com
REM Пароль приложения Gmail
set SMTP_PASS=babt hmnr vpxd qxyl
set FROM_EMAIL=kashkenov2006@gmail.com

echo Environment variables set successfully!
echo.
echo Note: For Gmail you need to:
echo 1. Enable 2-Step Verification in your Google Account
echo 2. Generate an App Password: Google Account -^> Security -^> App Passwords
echo 3. Use that App Password as SMTP_PASS
echo.
echo Current settings:
echo SMTP_HOST: %SMTP_HOST%
echo SMTP_PORT: %SMTP_PORT%
echo SMTP_USER: %SMTP_USER%
echo FROM_EMAIL: %FROM_EMAIL% 