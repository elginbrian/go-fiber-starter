package service

import (
	"errors"
	"fiber-starter/internal/domain"
	"fiber-starter/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(username, email, password string) error
    Login(email, password string) (string, error)
}

type authService struct {
    authRepo  repository.AuthRepository
    userRepo  repository.UserRepository
    jwtSecret string
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, jwtSecret string) AuthService {
    return &authService{authRepo: authRepo, userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(username, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &domain.User{
        Name:  username,
        Email: email,
        PasswordHash: string(hashedPassword),
    }

    return s.userRepo.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, error) {
    user, err := s.authRepo.GetUserByEmail(email) 
    if err != nil || user == nil {
        return "", errors.New("invalid email or password")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", errors.New("invalid email or password")
    }

    token := GenerateJWT(user.ID, s.jwtSecret)

    return token, nil
}
