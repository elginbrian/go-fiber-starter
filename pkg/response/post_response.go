package response

import "time"

type GetAllPostsResponse struct {
	Status string `json:"status" example:"success"`
	Data   []Post `json:"data"`
	Code   int    `json:"code" example:200`
}

type SearchPostsResponse struct {
	Status string `json:"status" example:"success"`
	Data   []Post `json:"data"`
	Code   int    `json:"code" example:200`
}

type GetPostByIDResponse struct {
	Status string `json:"status" example:"success"`
	Data   Post   `json:"data"`
	Code   int    `json:"code" example:200`
}

type CreatePostResponse struct {
	Status string `json:"status" example:"success"`
	Data   Post   `json:"data"`
	Code   int    `json:"code" example:201`
}

type UpdatePostResponse struct {
	Status string `json:"status" example:"success"`
	Data   Post   `json:"data"`
	Code   int    `json:"code" example:200`
}

type DeletePostResponse struct {
	Status string         `json:"status" example:"success"`
	Data   DeletePostData `json:"data"`
	Code   int            `json:"code" example:200`
}

type DeletePostData struct {
	Message string `json:"message" example:"Post deleted successfully"`
}

type Post struct {
	ID        string    `json:"id" example:"f9d6b52a-76a1-4b2b-9229-4c8db23a5ef2"`
	UserID    string    `json:"user_id" example:"2e0850c7-d213-4a91-9b78-bb86e3a6f0d3"`
	Caption   string    `json:"caption" example:"Had an amazing day at the beach!"`
	ImageURL  string    `json:"image_url" example:"https://example.com/images/beach.jpg"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-31T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-31T12:30:00Z"`
}
