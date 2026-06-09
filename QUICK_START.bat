@echo off
echo ====================================
echo Cinema Booking System - Local Setup
echo ====================================
echo.
echo Prerequisites:
echo - PostgreSQL running with 'cinema' database
echo - Node.js installed
echo - Go installed
echo.
echo Step 1: Start Backend
echo Open a NEW terminal and run:
echo   cd "E:\Project FE\Admin System\cinema-backend"
echo   go run cmd/api/main.go
echo.
echo Step 2: Start Frontend
echo Open ANOTHER NEW terminal and run:
echo   cd "E:\Project FE\Admin System\booking-room-admin"
echo   npm run dev
echo.
echo Step 3: Open in Browser
echo   http://localhost:3000
echo.
echo Test Credentials:
echo   Admin: admin@cinema.com / admin123
echo   Customer: Create new account
echo.
echo ====================================
pause
