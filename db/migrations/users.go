package migrations

const CreateUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
const InsertUsersData = `
INSERT INTO users (name, email, password_hash, created_at, updated_at)
VALUES 
    ('John Doe', 'john.doe@example.com', '$2a$10$D1g0OZLqH1rO5Gp2f9D5Fq2tqJb1h0B0VJhgV9AqE6qL6q8XgXy8G', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (name, email, password_hash, created_at, updated_at)
VALUES 
    ('Jane Smith', 'jane.smith@example.com', '$2a$10$E1g0OZLqH1rO5Gp2f9D5Fq2tqJb1h0B0VJhgV9AqE6qL6q8XgXy8G', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
`