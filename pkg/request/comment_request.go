package request

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required" example:"This is a great post!"`
}
