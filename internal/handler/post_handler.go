package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PostHandler handles post-related requests
type PostHandler struct {
	postService service.PostService
}

// NewPostHandler creates a new PostHandler instance
func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{postService: service}
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all posts from the database
// @Tags posts
// @Produce json
// @Success 200 {array} domain.Post "List of posts"
// @Failure 500 {object} fiber.Map{"error": "Internal server error"}
// @Router /api/posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postService.FetchAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Retrieves a post from the database by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} domain.Post "Post details"
// @Failure 400 {object} fiber.Map{"error": "Invalid post ID"}
// @Failure 404 {object} fiber.Map{"error": "Post not found"}
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}
	post, err := h.postService.FetchPostByID(postID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.JSON(post)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post in the database
// @Tags posts
// @Accept json
// @Produce json
// @Param post body domain.Post true "Post details"
// @Success 201 {object} domain.Post "Created post details"
// @Failure 400 {object} fiber.Map{"error": "Invalid input"}
// @Failure 500 {object} fiber.Map{"error": "Internal server error"}
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	createdPost, err := h.postService.CreatePost(post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdPost)
}

// UpdatePost godoc
// @Summary Update an existing post
// @Description Updates a post's details in the database by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body domain.Post true "Updated post details"
// @Success 200 {object} domain.Post "Updated post details"
// @Failure 400 {object} fiber.Map{"error": "Invalid post ID"}
// @Failure 404 {object} fiber.Map{"error": "Post not found"}
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}
	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	updatedPost, err := h.postService.UpdatePost(postID, post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.JSON(updatedPost)
}

// DeletePost godoc
// @Summary Delete a post by ID
// @Description Deletes a post from the database by its ID
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204 {object} string "No content"
// @Failure 400 {object} fiber.Map{"error": "Invalid post ID"}
// @Failure 404 {object} fiber.Map{"error": "Post not found"}
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}
	err = h.postService.DeletePost(postID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}