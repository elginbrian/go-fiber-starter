package handler

import (
	contract "fiber-starter/domain/contract"
	"fiber-starter/pkg/response"
	"fiber-starter/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type LikeHandler struct {
	likeService contract.ILikeService
	authService contract.IAuthService
}

func NewLikeHandler(likeService contract.ILikeService, authService contract.IAuthService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
		authService: authService,
	}
}

// LikePost godoc
// @Summary Like a post
// @Description Allows a user to like a post. Requires JWT authentication.
// @Tags likes
// @Param post_id path string true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.LikeResponse "Successfully liked post"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{post_id}/like [post]
func (h *LikeHandler) LikePost(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	postID := c.Params("post_id")
	if postID == "" {
		return response.ValidationError(c, "Post ID is required")
	}

	like, err := h.likeService.AddLike(user.ID, postID)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, like, fiber.StatusOK)
}

// UnlikePost godoc
// @Summary Unlike a post
// @Description Allows a user to remove their like from a post. Requires JWT authentication.
// @Tags likes
// @Param post_id path string true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.LikeResponse "Successfully unliked post"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{post_id}/unlike [post]
func (h *LikeHandler) UnlikePost(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	postID := c.Params("post_id")
	if postID == "" {
		return response.ValidationError(c, "Post ID is required")
	}

	err = h.likeService.RemoveLike(user.ID, postID)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, fiber.Map{"message": "Post unliked successfully"}, fiber.StatusOK)
}

// GetLikesByPostID godoc
// @Summary Get all likes for a post
// @Description Fetch all users who liked a specific post
// @Tags likes
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {array} response.GetAllLikesResponse "List of users who liked the post"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{post_id}/likes [get]
func (h *LikeHandler) GetLikesByPostID(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	if postID == "" {
		return response.ValidationError(c, "Post ID is required")
	}

	likes, err := h.likeService.GetLikesByPostID(postID)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, likes, fiber.StatusOK)
}

// GetLikesByUserID godoc
// @Summary Get all likes by a user
// @Description Fetch all posts liked by a specific user
// @Tags likes
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} response.GetAllLikesResponse "List of posts liked by the user"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/{user_id}/likes [get]
func (h *LikeHandler) GetLikesByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return response.ValidationError(c, "User ID is required")
	}

	likes, err := h.likeService.GetLikesByUserID(userID)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, likes, fiber.StatusOK)
}