-- Insert admin user
-- Password: admin123 (hashed with bcrypt)
INSERT INTO users (id, name, email, password, role, phone, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Admin',
    'admin@cinema.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'admin',
    '0123456789',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT DO NOTHING;

-- Insert test customer
-- Password: customer123
INSERT INTO users (id, name, email, password, role, phone, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'Test Customer',
    'customer@test.com',
    '$2a$10$OD1XtwQJgoocdY2YwrkffO7OGtGEi0cueQZSOMp6uzdSmqkkK47XC',
    'customer',
    '0987654321',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT DO NOTHING;
