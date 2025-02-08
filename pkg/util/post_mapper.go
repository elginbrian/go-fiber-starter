package util

import (
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/response"
)

func MapToPostResponse(posts []entity.Post) []response.Post {
	var postResponse []response.Post
	for _, post := range posts {
		postResponse = append(postResponse, response.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Caption:   post.Caption,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}
	return postResponse
}

func MapToSinglePostResponse(post entity.Post) response.Post {
	return response.Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Caption:   post.Caption,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
	}
}