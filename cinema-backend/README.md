# Cinema Booking Backend

Go backend API for cinema ticket booking system.

## Features

- Authentication (JWT)
- Role-based access (Admin/Customer)
- CRUD operations for Movies, Halls, Showtimes
- Ticket booking system
- Dashboard analytics
- Pagination, search, filtering

## Tech Stack

- Go 1.21+
- Gin Framework
- GORM (PostgreSQL)
- JWT Authentication
- CORS enabled

## Setup

1. Install Go 1.21 or later
2. Install PostgreSQL
3. Copy `.env.example` to `.env` and configure
4. Run:

```bash
go mod download
go run cmd/api/main.go
```

## API Endpoints

### Auth
- POST `/api/auth/register` - Register new customer
- POST `/api/auth/login` - Login
- GET `/api/auth/me` - Get current user (protected)
- POST `/api/auth/logout` - Logout (protected)

### Movies (Public)
- GET `/api/movies` - List all movies (with pagination, search, filter)
- GET `/api/movies/now-showing` - Now showing movies
- GET `/api/movies/coming-soon` - Coming soon movies
- GET `/api/movies/:id` - Get movie details

### Movies (Admin)
- POST `/api/movies` - Create movie
- PUT `/api/movies/:id` - Update movie
- DELETE `/api/movies/:id` - Delete movie

### Halls (Admin)
- GET `/api/halls` - List halls
- POST `/api/halls` - Create hall
- PUT `/api/halls/:id` - Update hall
- DELETE `/api/halls/:id` - Delete hall

### Showtimes
- GET `/api/showtimes` - List all (admin)
- GET `/api/showtimes/movie/:movieId` - By movie (customer)
- POST `/api/showtimes` - Create (admin)
- PUT `/api/showtimes/:id` - Update (admin)
- DELETE `/api/showtimes/:id` - Delete (admin)

### Tickets (Customer)
- POST `/api/tickets/book` - Book tickets
- GET `/api/tickets/my-tickets` - My tickets
- POST `/api/tickets/:id/cancel` - Cancel ticket

### Bookings
- GET `/api/bookings` - All bookings (admin)
- GET `/api/bookings/my-bookings` - My bookings (customer)
- POST `/api/bookings/:id/confirm` - Confirm (admin)
- POST `/api/bookings/:id/cancel` - Cancel
- POST `/api/bookings/:id/refund` - Refund (admin)

### Dashboard
- GET `/api/dashboard/admin-stats` - Admin stats
- GET `/api/dashboard/customer-stats` - Customer stats
- GET `/api/dashboard/top-movies` - Top movies
- GET `/api/dashboard/revenue-by-day` - Revenue chart data

## Database

PostgreSQL with tables:
- users
- movies
- cinema_halls
- showtimes
- tickets
- bookings
