package service

import (
	"context"
	"errors"
	"fiber-starter/internal/domain"
	"fiber-starter/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, email, password string) error
	Login(email, password string) (string, error)
	ChangePassword(userID int, oldPassword, newPassword string) error
}

type authService struct {
	authRepo  repository.AuthRepository
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{authRepo: authRepo, userRepo: userRepo, jwtSecret: jwtSecret}
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrIncorrectPassword  = errors.New("incorrect old password")
	ErrPasswordHashing    = errors.New("error hashing new password")
	ErrUpdatingUser       = errors.New("error updating user in database")
)

func (s *authService) Register(username, email, password string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = s.userRepo.CreateUser(ctx, user)
	return err
}

func (s *authService) Login(email, password string) (string, error) {
	ctx := context.Background()

	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ChangePassword(userID int, oldPassword, newPassword string) error {
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

func GenerateJWT(userID int, secret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}