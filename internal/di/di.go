package di

import (
	"fiber-starter/config"
	"fiber-starter/internal/handler"
	"fiber-starter/internal/repository"
	"fiber-starter/internal/service"
)

type Container struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewContainer() *Container {
	// Repositories
	userRepo := repository.NewUserRepository(config.DB)
	authRepo := repository.NewAuthRepository(config.DB)

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, config.GetEnv("JWT_SECRET", "secret"))

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
