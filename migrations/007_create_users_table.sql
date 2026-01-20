-- Create roles enum if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('admin', 'management', 'executive');
    END IF;
END $$;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    user_display_id VARCHAR(50) UNIQUE NOT NULL, -- This is the "User ID" provided by Admin
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    role user_role NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    password_set_token VARCHAR(255),
    password_set_token_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add user_id to expenses table
-- We don't make it NOT NULL yet because existing expenses don't have users
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id);

-- Seed default admin user
-- Password: password123 (hashed using bcrypt - $2a$10$7XvPkL.5O5.3.4.4.4.4.4) 
-- Wait, I'll just leave password_hash NULL and set a token for the default admin or a fixed hash.
-- Actually, the prompt says "ensure at least one Admin user is maintained in the database by default."
-- I'll create one admin with a default password for initial access, or better, set it as active.

INSERT INTO users (username, user_display_id, email, password_hash, role, is_active)
VALUES ('System Admin', 'admin01', 'admin@example.com', '$2a$10$Cmc1dCRV.yZHbV2Z0eHQ.uzQnBtY.Oeb0xa90n3gYW3MdOg9ERHM6', 'admin', TRUE)
ON CONFLICT (username) DO NOTHING;
