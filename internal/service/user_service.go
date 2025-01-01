package service

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/repository"
)

type UserService interface {
	FetchAllUsers() ([]domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) FetchAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAllUsers()
}
