package domain

import (
	"context"
	domain "fiber-starter/domain/entity"
)

type IUserRepository interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, id string, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, query string) ([]domain.User, error)
}

type IUserService interface {
	FetchAllUsers() ([]domain.User, error)
	FetchUserByID(id string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	UpdateUser(id string, user domain.User) (domain.User, error)
	DeleteUser(id string) error
	SearchUsers(query string) ([]domain.User, error)
}