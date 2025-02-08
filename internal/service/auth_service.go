package service

import (
	"context"
	"errors"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/util"

	"time"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo    contract.IUserRepository
	authRepo    contract.IAuthRepository
	jwtSecret   string
	refreshSecret string
}

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrPasswordHashing   = errors.New("error hashing password")
	ErrUpdatingUser      = errors.New("error updating user in database")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

func NewAuthService(userRepo contract.IUserRepository, authRepo contract.IAuthRepository, jwtSecret, refreshSecret string) contract.IAuthService {
	return &authService{userRepo: userRepo, authRepo: authRepo, jwtSecret: jwtSecret, refreshSecret: refreshSecret}
}

func (s *authService) Register(username, email, password string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashing
	}

	user := entity.User{
		Name:         username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = s.userRepo.CreateUser(ctx, user)
	return err
}

func (s *authService) Login(email, password string) (string, string, error) {
	ctx := context.Background()

	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrIncorrectPassword
	}

	accessToken, err := util.GenerateJWT(user.ID, s.jwtSecret, AccessTokenExpiration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.GenerateJWT(user.ID, s.refreshSecret, RefreshTokenExpiration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := util.ParseJWT(refreshToken, s.refreshSecret)
	if err != nil {
		return "", ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return util.GenerateJWT(userID, s.jwtSecret, AccessTokenExpiration)
}

func (s *authService) GetCurrentUser(ctx context.Context, token string) (*entity.User, error) {
	claims, err := util.ParseJWT(token, s.jwtSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) ChangePassword(userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrIncorrectPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashing
	}

	user.PasswordHash = string(hashedPassword)
	if _, err := s.userRepo.UpdateUser(ctx, user.ID, user); err != nil {
		return ErrUpdatingUser
	}

	return nil
}