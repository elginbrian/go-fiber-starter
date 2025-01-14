package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{postService: service}
}

type PostResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all posts from the database
// @Tags posts
// @Produce json
// @Success 200 {array} PostResponse "List of posts"
// @Failure 500 {object} ErrorResponse
// @Router /api/posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postService.FetchAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	var response []PostResponse
	for _, post := range posts {
		response = append(response, PostResponse{
			ID:      post.ID,
			Title:   post.Caption,
			Content: post.ImageURL,
		})
	}
	return c.JSON(response)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Retrieves a post from the database by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostResponse "Post details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid post ID"})
	}
	post, err := h.postService.FetchPostByID(postID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Post not found"})
	}
	response := PostResponse{
		ID:      post.ID,
		Title:   post.Caption,
		Content: post.ImageURL,
	}
	return c.JSON(response)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post in the database
// @Tags posts
// @Accept json
// @Produce json
// @Param post body domain.Post true "Post details"
// @Success 201 {object} PostResponse "Created post details"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid input"})
	}
	createdPost, err := h.postService.CreatePost(post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}
	response := PostResponse{
		ID:      createdPost.ID,
		Title:   createdPost.Caption,
		Content: createdPost.ImageURL,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdatePost godoc
// @Summary Update an existing post
// @Description Updates a post's details in the database by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body domain.Post true "Updated post details"
// @Success 200 {object} PostResponse "Updated post details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid post ID"})
	}
	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid input"})
	}
	updatedPost, err := h.postService.UpdatePost(postID, post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Post not found"})
	}
	response := PostResponse{
		ID:      updatedPost.ID,
		Title:   updatedPost.Caption,
		Content: updatedPost.ImageURL,
	}
	return c.JSON(response)
}

// DeletePost godoc
// @Summary Delete a post by ID
// @Description Deletes a post from the database by its ID
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204 {object} string "No content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid post ID"})
	}
	err = h.postService.DeletePost(postID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: "Post not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}