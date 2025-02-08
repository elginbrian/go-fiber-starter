package request

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
}
