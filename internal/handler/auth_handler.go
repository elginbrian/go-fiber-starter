package handler

import (
	"fiber-starter/internal/service"
	"fiber-starter/pkg/request"
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
// @Param request body request.UserRegistrationRequest true "User registration details"
// @Success 201 {object} response.RegisterResponse "Successful registration response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.UserRegistrationRequest

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

	return response.Success(c.Status(fiber.StatusCreated), response.RegisterData{
		Message: "User registered successfully",
	})
}

// @Summary Logs in an existing user
// @Description This endpoint allows a user to log in by providing their email and password. Upon successful login, a JWT token is generated and returned, which can be used for authenticated requests.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.UserLoginRequest true "User login details"
// @Success 201 {object} response.LoginResponse "Successful registration response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.UserLoginRequest

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

	return response.Success(c, response.LoginData{
		Token: "Bearer " + token,
	})
}

// @Summary Changes a user's password 
// @Description This endpoint allows a user to change their password by providing their old password, new password, and user ID. The request requires authentication via a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ChangePasswordRequest true "User change password details"
// @Security BearerAuth
// @Success 200 {object} response.ChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req request.ChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ValidationError(c, "Invalid request format")
	}

	validate := validator.New()
	if validationErrs := validate.Struct(req); validationErrs != nil {
		return response.ValidationError(c, validationErrs.Error())
	}

	if req.UserID == 0 {
		return response.ValidationError(c, "User ID is required")
	}

	userID := req.UserID

	if userID != c.Locals("user_id").(int) {
		return response.Error(c.Status(fiber.StatusUnauthorized), "Unauthorized to change this password")
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return response.Error(c, err.Error())
	}

	return response.Success(c, response.ChangePasswordData{
		Message: "Password changed successfully",
	})
}