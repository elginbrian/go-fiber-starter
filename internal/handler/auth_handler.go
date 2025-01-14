package handler

import (
	"fiber-starter/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Registers a new user
// @Description Registers a new user with a username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body struct { Username string `json:"username"`; Email string `json:"email"`; Password string `json:"password"` } true "User Registration Details"
// @Success 201 {object} fiber.Map{"message": "User registered successfully"}
// @Failure 400 {object} fiber.Map{"error": "Invalid request"}
// @Failure 500 {object} fiber.Map{"error": "Internal server error"}
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.authService.Register(req.Username, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login godoc
// @Summary Logs in a user
// @Description Logs in a user by verifying email and password, and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body struct { Email string `json:"email"`; Password string `json:"password"` } true "User Login Details"
// @Success 200 {object} fiber.Map{"token": "JWT token"}
// @Failure 400 {object} fiber.Map{"error": "Invalid request"}
// @Failure 401 {object} fiber.Map{"error": "Unauthorized"}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}