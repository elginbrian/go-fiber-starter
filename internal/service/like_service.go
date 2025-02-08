package service

import (
	"context"
	"errors"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
)

type likeService struct {
	likeRepo contract.ILikeRepository
}

func NewLikeService(repo contract.ILikeRepository) contract.ILikeService {
	return &likeService{likeRepo: repo}
}

func (s *likeService) GetLikesByPostID(postID string) ([]entity.Like, error) {
	ctx := context.Background()
	likes, err := s.likeRepo.GetLikesByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if len(likes) == 0 {
		return nil, errors.New("not found")
	}
	return likes, nil
}

func (s *likeService) GetLikesByUserID(userID string) ([]entity.Like, error) {
	ctx := context.Background()
	likes, err := s.likeRepo.GetLikesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(likes) == 0 {
		return nil, errors.New("not found")
	}
	return likes, nil
}

func (s *likeService) AddLike(userID, postID string) (entity.Like, error) {
	ctx := context.Background()
	like := entity.Like{
		UserID: userID,
		PostID: postID,
	}
	createdLike, err := s.likeRepo.AddLike(ctx, like)
	if err != nil {
		return entity.Like{}, err
	}
	return *createdLike, nil
}

func (s *likeService) RemoveLike(userID, postID string) error {
	ctx := context.Background()
	err := s.likeRepo.RemoveLike(ctx, userID, postID)
	if err != nil {
		return err
	}
	return nil
}