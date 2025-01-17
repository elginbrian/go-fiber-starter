package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"fiber-starter/pkg/response"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{postService: service}
}


// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all posts, including the user who created them, the caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Success 200 {object} map[string]interface{} "Successful fetch posts response" example({"status": "success", "data": [{"id": 1, "user_id": 1, "caption": "A beautiful view of the fjords", "image_url": "https://www.w3schools.com/w3images/fjords.jpg", "created_at": "2025-01-17 06:23:03", "updated_at": "2025-01-17 06:23:03"}, {"id": 2, "user_id": 2, "caption": "The city lights at night", "image_url": "https://www.w3schools.com/w3images/lights.jpg", "created_at": "2025-01-17 06:23:03", "updated_at": "2025-01-17 06:23:03"}]})
// @Router /api/posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postService.FetchAllPosts()
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	var postResponse []domain.PostResponse
	for _, post := range posts {
		postResponse = append(postResponse, domain.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Caption:   post.Caption,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(c, postResponse, fiber.StatusOK)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Retrieves a specific post by its ID, including its caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Successful fetch post response" example({"status": "success", "data": {"id": 1, "user_id": 1, "caption": "A beautiful view of the fjords", "image_url": "https://www.w3schools.com/w3images/fjords.jpg", "created_at": "2025-01-17 06:23:03", "updated_at": "2025-01-17 06:23:03"}})
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid post ID", fiber.StatusBadRequest)
	}
	post, err := h.postService.FetchPostByID(postID)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	postResponse := domain.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Caption:   post.Caption,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return response.Success(c, postResponse, fiber.StatusOK)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post with an optional image. The caption is required. If an image is provided, it will be uploaded to the server, and the URL will be returned in the response.
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param caption formData string true "Post caption"
// @Param image formData file false "Post image (optional)"
// @Security BearerAuth
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return response.ValidationError(c, "Invalid or missing user ID")
	}

	caption := c.FormValue("caption")
	if caption == "" {
		return response.ValidationError(c, "Caption is required")
	}

	var imageURL string

	file, err := c.FormFile("image")
	if err == nil {
		uploadDir := "./public/uploads/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				return response.Error(c, "Failed to create upload directory", fiber.StatusInternalServerError)
			}
		}

		sanitizedFileName := sanitizeFileName(file.Filename)

		savePath := uploadDir + sanitizedFileName
		if err := c.SaveFile(file, savePath); err != nil {
			return response.Error(c, "Failed to save image", fiber.StatusInternalServerError)
		}

		imageURL = "http://localhost:8084" + "/uploads/" + sanitizedFileName
	}

	post := domain.Post{
		UserID:   userID,
		Caption:  caption,
		ImageURL: imageURL, 
	}

	createdPost, err := h.postService.CreatePost(post)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	postResponse := domain.PostResponse{
		ID:        createdPost.ID,
		UserID:    createdPost.UserID,
		Caption:   createdPost.Caption,
		ImageURL:  createdPost.ImageURL,
		CreatedAt: createdPost.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: createdPost.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response.Success(c, postResponse, fiber.StatusCreated)
}

func sanitizeFileName(fileName string) string {
	sanitized := strings.ReplaceAll(fileName, " ", "_")
	sanitized = regexp.MustCompile(`[^a-zA-Z0-9\._-]`).ReplaceAllString(sanitized, "")
	return sanitized
}

// UpdatePost godoc
// @Summary Update an existing post
// @Description Updates the details of a post (caption and/or image). Only the creator of the post is allowed to update it.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body domain.Post true "Updated post details"
// @Security BearerAuth
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid post ID", fiber.StatusBadRequest)
	}

	var post domain.Post
	if err := c.BodyParser(&post); err != nil {
		return response.ValidationError(c, "Invalid input")
	}

	if post.UserID != userID {
		return response.Error(c, "Unauthorized to update this post", fiber.StatusUnauthorized)
	}

	updatedPost, err := h.postService.UpdatePost(postID, post)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	postResponse := domain.PostResponse{
		ID:        updatedPost.ID,
		UserID:    updatedPost.UserID,
		Caption:   updatedPost.Caption,
		ImageURL:  updatedPost.ImageURL,
		CreatedAt: updatedPost.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedPost.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return response.Success(c, postResponse)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes a post by its ID. Only the creator of the post is authorized to delete it.
// @Tags posts
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, "Invalid post ID", fiber.StatusBadRequest)
	}

	post, err := h.postService.FetchPostByID(postID)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	if post.UserID != userID {
		return response.Error(c, "Unauthorized to delete this post", fiber.StatusUnauthorized)
	}

	err = h.postService.DeletePost(postID)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}
	return response.Success(c, fiber.Map{"message": "Post deleted successfully"}, fiber.StatusNoContent)
}