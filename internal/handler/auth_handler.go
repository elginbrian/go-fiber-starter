package handler

import (
	"fiber-starter/internal/service"
	"fiber-starter/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type UserRegistrationRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req UserRegistrationRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ValidationError(c, "Invalid request format")
	}

	validate := validator.New()
	if validationErrs := validate.Struct(req); validationErrs != nil {
		return response.ValidationError(c, validationErrs.Error())
	}

	if err := h.authService.Register(req.Username, req.Email, req.Password); err != nil {
		return response.Error(c, err.Error())
	}

	return response.Success(c.Status(fiber.StatusCreated), UserRegistrationResponse{
		Message: "User registered successfully",
	})
}

// @Summary Logs in a user
// @Description Logs in a user by verifying email and password, and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body UserLoginRequest true "User Login Details"
// @Success 200 {object} UserLoginResponse
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ValidationError(c, "Invalid request format")
	}

	validate := validator.New()
	if validationErrs := validate.Struct(req); validationErrs != nil {
		return response.ValidationError(c, validationErrs.Error())
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c.Status(fiber.StatusUnauthorized), err.Error())
	}

	return response.Success(c, UserLoginResponse{
		Token: token,
	})
}