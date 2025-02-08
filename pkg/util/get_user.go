package util

import (
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func GetUserFromToken(c *fiber.Ctx, authService contract.IAuthService) (*entity.User, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) <= len("Bearer ") {
		return nil, response.Error(c, "Missing or invalid token", fiber.StatusUnauthorized)
	}
	token := authHeader[len("Bearer "):]
	ctx := c.Context()
	return authService.GetCurrentUser(ctx, token)
}
