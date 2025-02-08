package domain

import (
	"context"
	domain "fiber-starter/domain/entity"
)

type IPostRepository interface {
	FetchAllPosts(ctx context.Context) ([]domain.Post, error)
	FetchPostByID(ctx context.Context, postID string) (*domain.Post, error)
	FetchPostsByUserID(ctx context.Context, userID string) ([]domain.Post, error)
	CreatePost(ctx context.Context, post domain.Post) (*domain.Post, error)
	UpdatePost(ctx context.Context, postID string, post domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, postID string) error
	SearchPosts(ctx context.Context, query string) ([]domain.Post, error)
}

type IPostService interface {
	FetchAllPosts() ([]domain.Post, error)
	FetchPostByID(id string) (domain.Post, error)
	FetchPostsByUserID(userID string) ([]domain.Post, error)
	CreatePost(post domain.Post) (domain.Post, error)
	UpdatePost(id string, post domain.Post) (domain.Post, error)
	DeletePost(id string) error
	SearchPosts(query string) ([]domain.Post, error)
}