package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"fiber-starter/pkg/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves all users from the database
// @Tags users
// @Produce json
// @Success 200 {array} UserResponse "List of users with timestamps"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.FetchAllUsers()
	if err != nil {
		return response.Error(c, err.Error())
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Username:  user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return response.Success(c, userResponses)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieves a user from the database by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "User details with timestamps"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid user ID")
	}

	user, err := h.userService.FetchUserByID(userID)
	if err != nil {
		return response.Error(c, "User not found")
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Username:  user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response.Success(c, userResponse)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User details"
// @Success 201 {object} UserResponse "Created user details with timestamps"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return response.ValidationError(c, "Invalid input")
	}

	if validationErrs := response.ValidateStruct(user); validationErrs != nil {
		return response.ValidationError(c, validationErrs)
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		return response.Error(c, err.Error())
	}

	userResponse := UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	return response.Success(c.Status(fiber.StatusCreated), userResponse)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Updates a user's details in the database by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "Updated user details"
// @Success 200 {object} UserResponse "Updated user details with timestamps"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid user ID")
	}

	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return response.ValidationError(c, "Invalid input")
	}

	if validationErrs := response.ValidateStruct(user); validationErrs != nil {
		return response.ValidationError(c, validationErrs)
	}

	updatedUser, err := h.userService.UpdateUser(userID, user)
	if err != nil {
		return response.Error(c, "User not found")
	}

	userResponse := UserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return response.Success(c, userResponse)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Deletes a user from the database by their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 {string} string "No content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid user ID")
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		return response.Error(c, "User not found")
	}

	return response.Success(c, "User deleted successfully")
}