package domain

import "time"

type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Caption   string    `json:"caption,omitempty"`
    ImageURL  string    `json:"image_url,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}