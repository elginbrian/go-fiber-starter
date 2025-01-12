package di

import (
	"database/sql"
	"fiber-starter/internal/handler"
	"fiber-starter/internal/repository"
	"fiber-starter/internal/service"
)

type Container struct {
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

func NewContainer(db *sql.DB, jwtSecret string) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, jwtSecret)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}