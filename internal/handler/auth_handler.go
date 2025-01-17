package handler

import (
	"fiber-starter/internal/domain"
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

// @Summary Registers a new user
// @Description This endpoint allows users to create a new account by providing a username, email, and password. The registration data is validated, and upon successful registration, a success message is returned.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.UserRegistrationRequest true "User registration details"
// @Success 201 {object} map[string]interface{} "Successful registration response" example({"status": "success", "data": {"message": "User registered successfully"}})
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.UserRegistrationRequest

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

	return response.Success(c.Status(fiber.StatusCreated), domain.UserRegistrationResponse{
		Message: "User registered successfully",
	})
}

// @Summary Logs in an existing user
// @Description This endpoint allows a user to log in by providing their email and password. Upon successful login, a JWT token is generated and returned, which can be used for authenticated requests.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.UserLoginRequest true "User login details"
// @Success 201 {object} map[string]interface{} "Successful registration response" example({"status": "success", "data": {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzcxODE2NTAsInVzZXJfaWQiOjN9.nd-wN002UeYukRVwSY1jsNVAz9hcB2p24eHuvV5Pe6E"}})
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.UserLoginRequest

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

	return response.Success(c, domain.UserLoginResponse{
		Token: token,
	})
}