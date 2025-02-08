package domain

import (
	"context"
	domain "fiber-starter/domain/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type IAuthService interface {
	Register(username, email, password string) error
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	GetCurrentUser(ctx context.Context, token string) (*domain.User, error)
}