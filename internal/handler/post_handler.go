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

func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postService.FetchAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

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