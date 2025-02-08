package migrations

const CreateUsersTable = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL,  
    bio TEXT NOT NULL,                
    created_at TIMESTAMP,   
    updated_at TIMESTAMP    
);
`

const InsertUsersData = `
INSERT INTO users (name, email, password_hash, image_url, bio, created_at, updated_at)
VALUES 
    ('elginbrian', 'elginbrian49@gmail.com', '$2a$10$D1g0OZLqH1rO5Gp2f9D5Fq2tqJb1h0B0VJhgV9AqE6qL6q8XgXy8G', 
     'https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg', 
     'Hi there!', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (name, email, password_hash, image_url, bio, created_at, updated_at)
VALUES 
    ('midnightsparks', 'midnightsparks@example.com', '$2a$10$E1g0OZLqH1rO5Gp2f9D5Fq2tqJb1h0B0VJhgV9AqE6qL6q8XgXy8G', 
     'https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg', 
     'Hi there!', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
`
