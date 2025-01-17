package migrations

const CreatePostsTable = `
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT,
    caption TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, image_url, caption) 
);
`

const InsertPostsData = `
INSERT INTO posts (user_id, image_url, caption, created_at, updated_at)
VALUES 
    (1, 'https://wallpapers.com/images/featured/beautiful-scenery-wnxju2647uqrcccv.jpg', 'A beautiful view of the waterfall', NOW(), NOW())
ON CONFLICT (user_id, image_url, caption) DO NOTHING;

INSERT INTO posts (user_id, image_url, caption, created_at, updated_at)
VALUES 
    (2, 'https://cdn11.bigcommerce.com/s-nq6l4syi/images/stencil/1280x1280/products/97985/778893/147317-1024__68590.1674087138.jpg?c=2', 'Cool looking scenery', NOW(), NOW())
ON CONFLICT (user_id, image_url, caption) DO NOTHING;
`
