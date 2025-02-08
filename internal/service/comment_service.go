package service

import (
	"context"
	"errors"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
)

type commentService struct {
	commentRepo contract.ICommentRepository
}

func NewCommentService(repo contract.ICommentRepository) contract.ICommentService {
	return &commentService{commentRepo: repo}
}

func (s *commentService) GetCommentsByPostID(postID string) ([]entity.Comment, error) {
	ctx := context.Background()
	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("not found")
	}
	return comments, nil
}

func (s *commentService) GetCommentByID(commentID string) (entity.Comment, error) {
	ctx := context.Background()
	comment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return entity.Comment{}, err
	}
	if comment == nil {
		return entity.Comment{}, errors.New("not found")
	}
	return *comment, nil
}

func (s *commentService) CreateComment(comment entity.Comment) (entity.Comment, error) {
	ctx := context.Background()
	createdComment, err := s.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return entity.Comment{}, err
	}
	if createdComment == nil {
		return entity.Comment{}, errors.New("not found")
	}
	return *createdComment, nil
}

func (s *commentService) DeleteComment(commentID string) error {
	ctx := context.Background()
	err := s.commentRepo.DeleteComment(ctx, commentID)
	if err != nil {
		return err
	}
	return nil
}
