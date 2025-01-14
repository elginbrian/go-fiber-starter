package service

import (
	"context"
	"fiber-starter/internal/domain"
	"fiber-starter/internal/repository"
)

type UserService interface {
	FetchAllUsers() ([]domain.User, error)
	FetchUserByID(id int) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	UpdateUser(id int, user domain.User) (domain.User, error)
	DeleteUser(id int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) FetchAllUsers() ([]domain.User, error) {
	ctx := context.Background()
	return s.userRepo.GetAllUsers(ctx)
}

func (s *userService) FetchUserByID(id int) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) CreateUser(user domain.User) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) UpdateUser(id int, user domain.User) (domain.User, error) {
	ctx := context.Background()
	return s.userRepo.UpdateUser(ctx, id, user)
}

func (s *userService) DeleteUser(id int) error {
	ctx := context.Background()
	return s.userRepo.DeleteUser(ctx, id)
}
