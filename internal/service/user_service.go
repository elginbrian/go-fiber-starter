package service

import (
	"context"
	contract "fiber-starter/domain/contract"
	domain "fiber-starter/domain/entity"
)

type userService struct {
	userRepo contract.IUserRepository
}

func NewUserService(repo contract.IUserRepository) contract.IUserService {
	return &userService{userRepo: repo}
}

func (s *userService) FetchAllUsers() ([]domain.User, error) {
	ctx := context.Background()
	return s.userRepo.GetAllUsers(ctx)
}

func (s *userService) FetchUserByID(id string) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) CreateUser(user domain.User) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(id string, user domain.User) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.UpdateUser(ctx, id, user)
}

func (s *userService) DeleteUser(id string) error {
	ctx := context.Background()
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *userService) SearchUsers(query string) ([]domain.User, error) {
    ctx := context.Background()
    return s.userRepo.SearchUsers(ctx, query)
}
