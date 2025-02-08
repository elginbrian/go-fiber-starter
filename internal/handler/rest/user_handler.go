package handler

import (
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/response"
	"fiber-starter/pkg/util"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService contract.IUserService
	authService contract.IAuthService
}

func NewUserHandler(userService contract.IUserService, authService contract.IAuthService) *UserHandler {
    return &UserHandler{
        userService: userService,
        authService: authService,
    }
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users from the database.
// @Tags users
// @Produce json
// @Success 200 {object} response.GetAllUsersResponse "Successful fetch users response"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.FetchAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return response.Error(c, "Error fetching users", fiber.StatusInternalServerError)
	}

	var userResponses []response.User
	for _, user := range users {
		userResponses = append(userResponses, util.MapToUserResponse(user))
	}

	return response.Success(c, userResponses)
}

// GetUserByID godoc
// @Summary Get user details by ID
// @Description Retrieve the details of a specific user by their ID.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.GetUserByIDResponse "Successful fetch user by ID response"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.FetchUserByID(id)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)
		return response.Error(c, "User not found", fiber.StatusNotFound)
	}

	return response.Success(c, util.MapToUserResponse(user))
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update the bio, image_url, and/or username of the authenticated user. All fields are optional. If a field is not provided, the existing value will be retained.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param username formData string false "Updated username (optional)"
// @Param bio formData string false "Updated bio (optional)"
// @Param image formData file false "Updated image (optional)"
// @Security BearerAuth
// @Success 200 {object} response.UpdateUserResponse "Successful update user response"
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized or invalid token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	token, err := util.GetToken(c)
	if err != nil {
		return response.Error(c.Status(fiber.StatusUnauthorized), "Unauthorized: "+err.Error())
	}

	ctx := c.Context()
	user, err := h.authService.GetCurrentUser(ctx, token)
	if err != nil {
		log.Printf("Unauthorized access: %v", err)
		return response.Error(c.Status(fiber.StatusUnauthorized), "Unauthorized: "+err.Error())
	}

	username := c.FormValue("username")
	bio := c.FormValue("bio")
	imageFile, err := c.FormFile("image")
	if err != nil && imageFile != nil {
		log.Printf("Error parsing form data: %v", err)
		return response.ValidationError(c, "Invalid input")
	}

	if username != "" && (len(username) < 3 || len(username) > 50) {
		return response.ValidationError(c, "Username must be between 3 and 50 characters")
	}

	var imageURL string
	if imageFile != nil {
		imageURL, err = util.UploadProfileImage(c, user.ID, "./uploads/profiles/")
		if err != nil {
			return response.Error(c.Status(fiber.StatusInternalServerError), "Failed to upload image")
		}
	}

	updatedUser := entity.User{
		ID:        user.ID,
		Name:      username,
		Email:     user.Email,
		Bio:       bio,
		ImageURL:  imageURL,
		CreatedAt: user.CreatedAt,
	}
	
	if username == "" {
		updatedUser.Name = user.Name 
	}
	if bio == "" {
		updatedUser.Bio = user.Bio 
	}
	if imageFile == nil {
		updatedUser.ImageURL = user.ImageURL 
	}

	updatedUser, err = h.userService.UpdateUser(user.ID, updatedUser)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return response.Error(c.Status(fiber.StatusInternalServerError), fmt.Sprintf("Error updating user: %v", err))
	}

	return response.Success(c, util.MapToUserResponse(updatedUser))
}

// SearchUsers godoc
// @Summary Search users
// @Description Search for users by their name or email.
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} response.SearchUsersResponse "Successful search response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /search/users [get]
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return response.Error(c, "Query parameter is required", fiber.StatusBadRequest)
	}

	log.Printf("Received search query: %s", query)

	users, err := h.userService.SearchUsers(query)
	if err != nil {
		if err.Error() == "no users found" {
			return response.Error(c, "No users found for the given search query", fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	var userResponses []response.User
	for _, user := range users {
		userResponses = append(userResponses, util.MapToUserResponse(user))
	}

	return response.Success(c, userResponses)
}