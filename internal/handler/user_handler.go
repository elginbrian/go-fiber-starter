package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves all users from the database
// @Tags users
// @Produce json
// @Success 200 {array} domain.User "List of users"
// @Failure 500 {object} fiber.Map{"error": "Internal server error"}
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.FetchAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Retrieves a user from the database by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User "User details"
// @Failure 400 {object} fiber.Map{"error": "Invalid user ID"}
// @Failure 404 {object} fiber.Map{"error": "User not found"}
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	user, err := h.userService.FetchUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User details"
// @Success 201 {object} domain.User "Created user details"
// @Failure 400 {object} fiber.Map{"error": "Invalid input"}
// @Failure 500 {object} fiber.Map{"error": "Internal server error"}
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Updates a user's details in the database by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body domain.User true "Updated user details"
// @Success 200 {object} domain.User "Updated user details"
// @Failure 400 {object} fiber.Map{"error": "Invalid user ID"}
// @Failure 404 {object} fiber.Map{"error": "User not found"}
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	updatedUser, err := h.userService.UpdateUser(userID, user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(updatedUser)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Deletes a user from the database by their ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 {object} string "No content"
// @Failure 400 {object} fiber.Map{"error": "Invalid user ID"}
// @Failure 404 {object} fiber.Map{"error": "User not found"}
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	err = h.userService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
