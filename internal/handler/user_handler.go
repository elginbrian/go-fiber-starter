package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"fiber-starter/pkg/request"
	"fiber-starter/pkg/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func parseUserID(c *fiber.Ctx) (string, error) {
	id := c.Params("id")
	return id, nil
}

// GetAllUsers godoc
// @Summary Retrieve a list of all users
// @Description Fetches all user records from the database and returns them with timestamps for creation and update.
// @Tags users
// @Produce json
// @Success 200 {object} response.GetAllUsersResponse "Successful fetch users response" 
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.FetchAllUsers()
	if err != nil {
		return response.Error(c, err.Error())
	}

	var userResponses []domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, domain.UserResponse{
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
// @Param id path string true "User ID"
// @Success 200 {object} response.GetUserByIDResponse "Successful fetch user by ID response" 
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.FetchUserByID(id)
	if err != nil {
		return response.Error(c, "User not found")
	}

	userResponse := domain.UserResponse{
		ID:        user.ID,
		Username:  user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response.Success(c, userResponse)
}

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

	userResponse := domain.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	return response.Success(c.Status(fiber.StatusCreated), userResponse)
}

// UpdateUser godoc
// @Summary Update an existing user's username
// @Description Updates the user's username. Users can only modify their own account.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "Request body with updated username"
// @Security BearerAuth
// @Success 200 {object} response.UpdateUserResponse "Successful update user response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    userID, err := parseUserID(c)
    if err != nil {
        return response.Error(c, err.Error(), fiber.StatusBadRequest)
    }

    authenticatedUserID, ok := c.Locals("user_id").(string)
    if !ok {
        return response.Error(c, "Unauthorized access", fiber.StatusUnauthorized)
    }

    if authenticatedUserID != userID {
        return response.Error(c, "You are not authorized to update this user", fiber.StatusForbidden)
    }

    existingUser, err := h.userService.FetchUserByID(userID)
    if err != nil {
        return response.Error(c, fmt.Sprintf("Error fetching user: %v", err), fiber.StatusInternalServerError)
    }

	if existingUser.ID == "" {
        return response.Error(c, "User not found", fiber.StatusNotFound)
    }

    var payload request.UpdateUserRequest

    if err := c.BodyParser(&payload); err != nil {
        return response.ValidationError(c, "Invalid input, expected JSON with 'username'")
    }

    if len(payload.Username) < 3 || len(payload.Username) > 50 {
        return response.ValidationError(c, "Username must be between 3 and 50 characters")
    }

    updatedUser := domain.User{
        ID:   userID,
        Name: payload.Username,
		Email: existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
    }

    updatedUser, err = h.userService.UpdateUser(userID, updatedUser)
    if err != nil {
        return response.Error(c, fmt.Sprintf("Error updating user: %v", err), fiber.StatusInternalServerError)
    }

    userResponse := domain.UserResponse{
        ID:        updatedUser.ID,
        Username:  updatedUser.Name,
        Email:     updatedUser.Email,
        CreatedAt: updatedUser.CreatedAt,
        UpdatedAt: updatedUser.UpdatedAt,
    }

    return response.Success(c, userResponse)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusBadRequest)
	}

	authenticatedUserID := c.Locals("userID").(string)

	if authenticatedUserID != userID {
		return response.Error(c, "You are not authorized to delete this user", fiber.StatusForbidden)
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		return response.Error(c, "User not found", fiber.StatusNotFound)
	}

	return response.Success(c, "User deleted successfully")
}

// SearchUsers godoc
// @Summary Search for users by name or email
// @Description Retrieves users that match the given search query in name or email.
// @Tags users
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} domain.UserResponse "Successful search response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/search/users [get]
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
    query := c.Query("query")

    if query == "" {
        return response.Error(c, "Query parameter is required", fiber.StatusBadRequest)
    }
    
    fmt.Println("Received search query:", query)

    users, err := h.userService.SearchUsers(query)
    if err != nil {
        if err.Error() == "no users found" {
            return response.Error(c, "No users found for the given search query", fiber.StatusNotFound)
        }
        return response.Error(c, err.Error(), fiber.StatusInternalServerError)
    }

    var userResponses []domain.UserResponse
    for _, user := range users {
        userResponses = append(userResponses, domain.UserResponse{
            ID:        user.ID,
            Username:  user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        })
    }

    return response.Success(c, userResponses)
}
