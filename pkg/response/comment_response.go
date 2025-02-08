package response

import "time"

type GetCommentsResponse struct {
	Status string    `json:"status" example:"success"`
	Data   []Comment `json:"data"`
	Code   int       `json:"code" example:200`
}

type CreateCommentResponse struct {
	Status string  `json:"status" example:"success"`
	Data   Comment `json:"data"`
	Code   int     `json:"code" example:201`
}

type DeleteCommentResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Comment deleted successfully"`
	Code    int    `json:"code" example:200`
}

type Comment struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string    `json:"user_id" example:"b3d1a42b-6871-4a47-bec3-6df0980a9c75"`
	PostID    string    `json:"post_id" example:"c6f7c988-233f-4f3c-a74d-17f72e4a1b56"`
	Content   string    `json:"content" example:"This is a comment!"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-31T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-01-31T12:30:00Z"`
}