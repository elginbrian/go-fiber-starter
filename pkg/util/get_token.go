package util

import (
	"fiber-starter/pkg/response"

	fiber "github.com/gofiber/fiber/v2"
)

func GetToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) <= len("Bearer ") {
		return "", response.Error(c, "Missing or invalid token", fiber.StatusUnauthorized)
	}
	return authHeader[len("Bearer "):], nil
}

