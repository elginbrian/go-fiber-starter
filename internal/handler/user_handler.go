package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves all users from the database
// @Tags users
// @Produce json
// @Success 200 {array} UserResponse "List of users"
// @Failure 500 {object} ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.FetchAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:       user.ID,
			Username: user.Name,
			Email:    user.Email,
		})
	}
	return c.JSON(response)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieves a user from the database by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse "User details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid user ID"})
	}
	user, err := h.userService.FetchUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "User not found"})
	}
	response := UserResponse{
		ID:       user.ID,
		Username: user.Name,
		Email:    user.Email,
	}
	return c.JSON(response)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User details"
// @Success 201 {object} UserResponse "Created user details"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid input"})
	}
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	response := UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Name,
		Email:    createdUser.Email,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Updates a user's details in the database by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "Updated user details"
// @Success 200 {object} UserResponse "Updated user details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid user ID"})
	}
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid input"})
	}
	updatedUser, err := h.userService.UpdateUser(userID, user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "User not found"})
	}
	response := UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Name,
		Email:    updatedUser.Email,
	}
	return c.JSON(response)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Deletes a user from the database by their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 {object} string "No content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid user ID"})
	}
	err = h.userService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "User not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}