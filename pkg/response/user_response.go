package response

import "time"

type GetAllUsersResponse struct {
	Status string `json:"status" example:"success"`
	Data   []User `json:"data"`
	Code   int    `json:"code" example:200`
}

type SearchUsersResponse struct {
	Status string `json:"status" example:"success"`
	Data   []User `json:"data"`
	Code   int    `json:"code" example:200`
}

type GetUserByIDResponse struct {
	Status string `json:"status" example:"success"`
	Data   User   `json:"data"`
	Code   int    `json:"code" example:200`
}

type UpdateUserResponse struct {
	Status string `json:"status" example:"success"`
	Data   User   `json:"data"`
	Code   int    `json:"code" example:200`
}

type DeleteUserResponse struct {
	Status string `json:"status" example:"success"`
	Data   string `json:"data" example:"User deleted successfully"`
	Code   int    `json:"code" example:200`
}

type User struct { 
	ID        string     `json:"id" example:"3d5a8b92-f1c5-4dbe-a2a7-1d9a8c743e9b"`
	Username  string     `json:"username" example:"john_doe"`
	Email     string     `json:"email" example:"john.doe@example.com"`
	ImageURL  string     `json:"image_url" example:"https://example.com/profile.jpg"`
	Bio       string     `json:"bio" example:"Hi there!"`
	CreatedAt time.Time  `json:"created_at" example:"2025-01-31T12:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2025-01-31T12:30:00Z"`
}