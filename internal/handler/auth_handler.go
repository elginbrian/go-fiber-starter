package handler

import (
	"fiber-starter/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	Message string `json:"message"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Registers a new user
// @Description Registers a new user with a username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body UserRegistrationRequest true "User Registration Details"
// @Success 201 {object} UserRegistrationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req UserRegistrationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request"})
	}

	if err := h.authService.Register(req.Username, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(UserRegistrationResponse{Message: "User registered successfully"})
}

// @Summary Logs in a user
// @Description Logs in a user by verifying email and password, and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body UserLoginRequest true "User Login Details"
// @Success 200 {object} UserLoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request"})
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.JSON(UserLoginResponse{Token: token})
}