package service

import (
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
	return s.userRepo.GetAllUsers()
}

func (s *userService) FetchUserByID(id int) (domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) CreateUser(user domain.User) (domain.User, error) {
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(id int, user domain.User) (domain.User, error) {
	return s.userRepo.UpdateUser(id, user)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}
