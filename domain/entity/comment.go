package domain

import "time"

type Comment struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    PostID    string    `json:"post_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}