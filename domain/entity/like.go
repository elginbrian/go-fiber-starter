package domain

import "time"

type Like struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    PostID    string    `json:"post_id"`
    CreatedAt time.Time `json:"created_at"`
}
