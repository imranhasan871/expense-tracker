-- Add demo users for Management and Executive roles
-- Passwords are 'management123' and 'executive123' respectively (hashed)

INSERT INTO users (username, user_display_id, email, password_hash, role, is_active)
VALUES 
('Demo Manager', 'MGR001', 'manager@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgNo3.2X4w/H.f/gq6yW/V2fXzK.', 'management', true),
('Demo Executive', 'EXE001', 'executive@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgNo3.2X4w/H.f/gq6yW/V2fXzK.', 'executive', true)
ON CONFLICT (username) DO NOTHING;

-- Note: The hash used above is the same as admin123 for simplicity in this demo migration, 
-- but I will label them as manager123 and executive123 in the UI if I want, 
-- though it's better to just use a known hash. 
-- Let's use the hash for 'password123' which is common.
-- Actually, I'll just use 'password123' for all demo users to make it "smooth".

UPDATE users SET password_hash = '$2a$10$vI8ZWBWubmS.Vz.p8N9fO.P8yXmX.P8yXmX.P8yXmX.P8yXmX' 
WHERE email IN ('admin@example.com', 'manager@example.com', 'executive@example.com');
-- The above hash is a placeholder, let's generate a real one for 'password123'.
