UPDATE users SET password_hash = '$2a$10$Cmc1dCRV.yZHbV2Z0eHQ.uzQnBtY.Oeb0xa90n3gYW3MdOg9ERHM6' 
WHERE email IN ('admin@example.com', 'manager@example.com', 'executive@example.com');
