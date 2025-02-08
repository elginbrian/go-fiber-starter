package response

import "time"

type LikeResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Like added successfully"`
	Code    int    `json:"code" example:200`
}

type GetAllLikesResponse struct {
	Status string `json:"status" example:"success"`
	Data   []Like `json:"data"`
	Code   int    `json:"code" example:200`
}

type Like struct {
	ID        string    `json:"id" example:"a1f5e4b3-8d2a-4c39-91a2-47b36295d8a3"`
	UserID    string    `json:"user_id" example:"b3d1a42b-6871-4a47-bec3-6df0980a9c75"`
	PostID    string    `json:"post_id" example:"c6f7c988-233f-4f3c-a74d-17f72e4a1b56"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-31T12:00:00Z"`
}
