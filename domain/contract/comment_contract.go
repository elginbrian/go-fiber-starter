package domain

import (
	"context"
	domain "fiber-starter/domain/entity"
)

type ICommentRepository interface {
	GetCommentsByPostID(ctx context.Context, postID string) ([]domain.Comment, error)
	GetCommentByID(ctx context.Context, commentID string) (*domain.Comment, error)
	CreateComment(ctx context.Context, comment domain.Comment) (*domain.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
}


type ICommentService interface {
	GetCommentsByPostID(postID string) ([]domain.Comment, error)
	GetCommentByID(commentID string) (domain.Comment, error)
	CreateComment(comment domain.Comment) (domain.Comment, error)
	DeleteComment(commentID string) error
}