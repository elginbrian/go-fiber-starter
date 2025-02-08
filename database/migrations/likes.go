package migrations

const CreateLikesTable = `
    CREATE EXTENSION IF NOT EXISTS pgcrypto; 

	CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,  
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, post_id)
);
`