package request

type UpdatePostRequest struct {
	Caption string `json:"caption" validate:"required,min=1" example:"Had an amazing trip to the mountains!"`
}
