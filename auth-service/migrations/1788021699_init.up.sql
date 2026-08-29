CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    full_name VARCHAR(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);