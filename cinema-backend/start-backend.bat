@echo off
setlocal

:: Set database URL
set DATABASE_URL=postgres://postgres:Tuananh6803%%40@localhost:5432/cinema?sslmode=disable

:: Kill any existing cinema-backend process
taskkill /F /IM cinema-backend.exe 2>nul

:: Wait a moment
timeout /t 2 /nobreak >nul

:: Start backend
cd /d "E:\Project FE\Admin System\cinema-backend"
start /B cinema-backend.exe

:: Wait and check
timeout /t 3 /nobreak >nul
echo Backend started on port 3001

:: Keep window open
echo.
echo Press any key to stop backend...
pause >nul

:: Kill backend when user presses key
taskkill /F /IM cinema-backend.exe 2>nul
echo Backend stopped.