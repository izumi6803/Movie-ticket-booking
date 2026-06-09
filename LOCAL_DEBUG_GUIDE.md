# Local Development Setup Guide

## Prerequisites
- Node.js 18+ (for Frontend)
- Go 1.21+ (for Backend)
- PostgreSQL 12+ (for Database)

## Step 1: Set up PostgreSQL Database

### Option A: Using Docker (Recommended)
```bash
docker run --name cinema-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=cinema \
  -p 5432:5432 \
  -d postgres:15
```

### Option B: Using Local PostgreSQL
Create a database named `cinema`:
```bash
createdb -U postgres cinema
```

## Step 2: Start the Backend (Go)

Open Terminal 1:
```bash
cd "E:\Project FE\Admin System\cinema-backend"
go mod download
go run cmd/api/main.go
```

The backend should start on: **http://localhost:3001**

You should see output like:
```
[GIN-debug] Listening and serving HTTP on :3001
```

### Verify Backend is Running:
```bash
curl http://localhost:3001/api/health
```

Expected response:
```json
{"status":"ok"}
```

## Step 3: Start the Frontend (Next.js)

Open Terminal 2:
```bash
cd "E:\Project FE\Admin System\booking-room-admin"
npm run dev
```

The frontend should start on: **http://localhost:3000**

You should see:
```
▲ Next.js 16.2.6
- Local:        http://localhost:3000
```

## Step 4: Open in Browser

Open: **http://localhost:3000**

## Testing the Application

### Test Credentials (if seeded):
- **Admin:**
  - Email: `admin@cinema.com`
  - Password: `admin123`

- **Customer:** (Create a new account or use existing)

### Test Payment Flow:
1. Login as customer
2. Book a ticket
3. Go to checkout/payment
4. Simulate payment success
5. Check for UI bugs and console errors

## Debugging Tools

### Frontend Debugging:
1. Open DevTools: **F12**
2. Go to **Console** tab to see errors
3. Go to **Network** tab to see API calls
4. Go to **Elements** tab to inspect CSS issues

### Backend Debugging:
Check the terminal where backend is running for:
- Request logs
- Database errors
- API response logs

## Environment Variables

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:3001/api
NEXT_PUBLIC_WS_URL=ws://localhost:3001
```

### Backend (.env)
```
PORT=3001
DATABASE_URL=postgres://postgres:postgres@localhost:5432/cinema?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

## Common Issues

### Backend won't start
- Check if PostgreSQL is running
- Check if port 3001 is already in use
- Run `go mod download` to ensure dependencies

### Frontend can't connect to backend
- Make sure backend is running on port 3001
- Check `.env.local` has correct `NEXT_PUBLIC_API_URL`
- Check browser console for CORS errors

### Database errors
- Make sure PostgreSQL is running
- Verify DATABASE_URL in .env is correct
- Check database exists: `psql -U postgres -l`

## Stopping Services

Press **Ctrl+C** in each terminal to stop the services.

---

Happy debugging! 🐛🔧
