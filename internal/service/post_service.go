package service

import (
	"context"
	"errors"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
)

type postService struct {
	postRepo contract.IPostRepository
}

func NewPostService(repo contract.IPostRepository) contract.IPostService {
	return &postService{postRepo: repo}
}

func (s *postService) FetchAllPosts() ([]entity.Post, error) {
	ctx := context.Background()
	return s.postRepo.FetchAllPosts(ctx)
}

func (s *postService) FetchPostByID(id string) (entity.Post, error) {
	ctx := context.Background()
	post, err := s.postRepo.FetchPostByID(ctx, id)
	if err != nil {
		return entity.Post{}, err
	}
	if post == nil {
		return entity.Post{}, errors.New("not found")
	}
	return *post, nil
}

func (s *postService) FetchPostsByUserID(userID string) ([]entity.Post, error) {
	ctx := context.Background()
	posts, err := s.postRepo.FetchPostsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("not found")
	}
	return posts, nil
}

func (s *postService) CreatePost(post entity.Post) (entity.Post, error) {
	ctx := context.Background()
	createdPost, err := s.postRepo.CreatePost(ctx, post)
	if err != nil {
		return entity.Post{}, err
	}
	if createdPost == nil {
		return entity.Post{}, errors.New("not found")
	}
	return *createdPost, nil
}

func (s *postService) UpdatePost(id string, post entity.Post) (entity.Post, error) {
	ctx := context.Background()
	updatedPost, err := s.postRepo.UpdatePost(ctx, id, post)
	if err != nil {
		return entity.Post{}, err
	}
	if updatedPost == nil {
		return entity.Post{}, errors.New("not found")
	}
	return *updatedPost, nil
}

func (s *postService) DeletePost(id string) error {
	ctx := context.Background()
	err := s.postRepo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) SearchPosts(query string) ([]entity.Post, error) {
	ctx := context.Background()
	posts, err := s.postRepo.SearchPosts(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("not found")
	}
	return posts, nil
}