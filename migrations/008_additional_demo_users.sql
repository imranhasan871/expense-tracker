-- Add demo users for Management and Executive roles
-- Passwords are 'password123' for all demo users (hashed)

INSERT INTO users (username, user_display_id, email, password_hash, role, is_active)
VALUES 
('Demo Manager', 'MGR001', 'manager@example.com', '$2a$10$Cmc1dCRV.yZHbV2Z0eHQ.uzQnBtY.Oeb0xa90n3gYW3MdOg9ERHM6', 'management', true),
('Demo Executive', 'EXE001', 'executive@example.com', '$2a$10$Cmc1dCRV.yZHbV2Z0eHQ.uzQnBtY.Oeb0xa90n3gYW3MdOg9ERHM6', 'executive', true)
ON CONFLICT (username) DO NOTHING;

UPDATE users SET password_hash = '$2a$10$Cmc1dCRV.yZHbV2Z0eHQ.uzQnBtY.Oeb0xa90n3gYW3MdOg9ERHM6' 
WHERE email IN ('admin@example.com', 'manager@example.com', 'executive@example.com');
-- The above hash is for 'password123'.
