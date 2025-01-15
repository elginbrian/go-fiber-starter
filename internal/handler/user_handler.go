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
// @Summary Retrieve a list of all users
// @Description Fetches all user records from the database and returns them with timestamps for creation and update.
// @Tags users
// @Produce json
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
// @Summary Retrieve user details by user ID
// @Description Fetches a specific user record from the database by the provided ID and returns the user's details, including timestamps.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
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
// @Summary Create a new user record
// @Description Accepts user details from the request body, validates the input, and creates a new user in the database. Returns the created user's details along with timestamps.
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User details"
// @Security BearerAuth
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
// @Summary Update an existing user's details
// @Description Accepts the user ID as a path parameter, along with updated user data in the request body, and updates the user record in the database. Returns the updated user's details with timestamps.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "Updated user details"
// @Security BearerAuth
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
// @Summary Delete a user record by ID
// @Description Accepts the user ID as a path parameter and deletes the user record from the database. Returns a success message upon deletion.
// @Tags users
// @Param id path int true "User ID"
// @Security BearerAuth
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