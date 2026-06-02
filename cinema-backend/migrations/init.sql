-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'customer',
    phone VARCHAR(50),
    avatar VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Movies table
CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL,
    genre TEXT[],
    rating VARCHAR(10),
    poster_url VARCHAR(500),
    trailer_url VARCHAR(500),
    release_date DATE,
    director VARCHAR(255),
    movie_cast TEXT[],
    status VARCHAR(50) DEFAULT 'coming_soon',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Theaters table
CREATE TABLE IF NOT EXISTS theaters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    total_screens INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Screens table
CREATE TABLE IF NOT EXISTS screens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    theater_id UUID REFERENCES theaters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    screen_type VARCHAR(50) DEFAULT 'standard',
    total_rows INTEGER NOT NULL,
    seats_per_row INTEGER NOT NULL,
    total_seats INTEGER NOT NULL,
    sound_system VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seats table
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    screen_id UUID REFERENCES screens(id) ON DELETE CASCADE,
    row_label VARCHAR(10) NOT NULL,
    seat_number INTEGER NOT NULL,
    seat_type VARCHAR(50) DEFAULT 'standard',
    price_multiplier DECIMAL(3,2) DEFAULT 1.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(screen_id, row_label, seat_number)
);

-- Showtimes table
CREATE TABLE IF NOT EXISTS showtimes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
    screen_id UUID REFERENCES screens(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    base_ticket_price DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Concessions table
CREATE TABLE IF NOT EXISTS concessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    image_url VARCHAR(500),
    stock_quantity INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    showtime_id UUID REFERENCES showtimes(id) ON DELETE CASCADE,
    booking_code VARCHAR(50) UNIQUE NOT NULL DEFAULT 'BOOK' || SUBSTRING(MD5(RANDOM()::TEXT), 1, 6),
    total_ticket_price DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_concession_price DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'pending',
    payment_status VARCHAR(50) DEFAULT 'pending',
    qr_code TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Booking seats
CREATE TABLE IF NOT EXISTS booking_seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id UUID REFERENCES seats(id) ON DELETE CASCADE,
    seat_label VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    UNIQUE(booking_id, seat_id)
);

-- Order items
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    concession_id UUID REFERENCES concessions(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
