package migrations

const CreatePostsTable = `
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT,
    caption TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
`

const InsertPostsData = `
INSERT INTO IF NOT EXISTS posts (user_id, image_url, caption, created_at, updated_at)
VALUES 
    (1, 'https://www.w3schools.com/w3images/fjords.jpg', 'A beautiful view of the fjords', NOW(), NOW());

INSERT INTO IF NOT EXISTS posts (user_id, image_url, caption, created_at, updated_at)
VALUES 
    (2, 'https://www.w3schools.com/w3images/lights.jpg', 'The city lights at night', NOW(), NOW());
`
