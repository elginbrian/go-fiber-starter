package response

import "time"

type GetAllPostsResponse struct {
	Status string `json:"status"`
	Data   []Post `json:"data"`
}

type GetPostByIDResponse struct {
	Status string `json:"status"`
	Data   Post   `json:"data"`
}

type CreatePostResponse struct {
	Status string `json:"status"`
	Data   Post   `json:"data"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Caption   string    `json:"caption"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}