CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,                    
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    image_url TEXT,                           
    caption TEXT,                             
    created_at TIMESTAMP DEFAULT NOW(),       
    updated_at TIMESTAMP DEFAULT NOW()        
);
