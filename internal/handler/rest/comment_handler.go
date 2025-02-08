package handler

import (
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/request"
	"fiber-starter/pkg/response"
	"fiber-starter/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	commentService contract.ICommentService
	authService    contract.IAuthService
}

func NewCommentHandler(commentService contract.ICommentService, authService contract.IAuthService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		authService:    authService,
	}
}

// GetCommentsByPostID godoc
// @Summary Get comments for a post
// @Description Retrieve all comments related to a specific post.
// @Tags comments
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {array} response.GetCommentsResponse "List of comments"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	comments, err := h.commentService.GetCommentsByPostID(postID)
	if err != nil {
		return response.Error(c, "No comments found", fiber.StatusNotFound)
	}

	return response.Success(c, comments, fiber.StatusOK)
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a comment for a post. Requires authentication.
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param request body request.CreateCommentRequest true "Comment request body"
// @Security BearerAuth
// @Success 201 {object} response.CreateCommentResponse "Created comment response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{post_id}/comments [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	var input request.CreateCommentRequest
	if err := c.BodyParser(&input); err != nil {
		return response.HandleValidationError(c, "Invalid input")
	}

	if input.Content == "" {
		return response.HandleValidationError(c, "Comment content cannot be empty")
	}

	comment := entity.Comment{
		PostID:  c.Params("post_id"),
		UserID:  user.ID,
		Content: input.Content,
	}

	createdComment, err := h.commentService.CreateComment(comment)
	if err != nil {
		return response.Error(c, "Failed to create comment", fiber.StatusInternalServerError)
	}

	return response.Success(c, createdComment, fiber.StatusCreated)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment by its ID. Only the comment creator can delete it. Requires authentication.
// @Tags comments
// @Param id path string true "Comment ID"
// @Security BearerAuth
// @Success 204 {object} response.DeleteCommentResponse "Successful deletion response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	id := c.Params("id")
	comment, err := h.commentService.GetCommentByID(id)
	if err != nil {
		return response.Error(c, "Comment not found", fiber.StatusNotFound)
	}

	if comment.UserID != user.ID {
		return response.Error(c, "Unauthorized to delete this comment", fiber.StatusUnauthorized)
	}

	if err := h.commentService.DeleteComment(id); err != nil {
		return response.Error(c, "Failed to delete comment", fiber.StatusInternalServerError)
	}

	return response.Success(c, fiber.Map{"message": "Comment deleted successfully"}, fiber.StatusNoContent)
}