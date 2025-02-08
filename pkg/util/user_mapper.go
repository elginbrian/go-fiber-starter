package util

import (
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/response"
)

func MapToUserResponse(user entity.User) response.User {
	return response.User{
		ID:        user.ID,
		Username:  user.Name,
		Email:     user.Email,
		Bio:       user.Bio,
		ImageURL:  user.ImageURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}