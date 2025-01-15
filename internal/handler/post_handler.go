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

type PostResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Caption   string `json:"caption,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
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
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	var postResponse []PostResponse
	for _, post := range posts {
		postResponse = append(postResponse, PostResponse{
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
		return response.Error(c, "Invalid post ID", fiber.StatusBadRequest)
	}
	post, err := h.postService.FetchPostByID(postID)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	postResponse := PostResponse{
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
// @Description Creates a new post in the database
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param caption formData string true "Caption"
// @Param image formData file true "Post image"  
// @Success 201 {object} PostResponse "Created post details"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

	postResponse := PostResponse{
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

	postResponse := PostResponse{
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
// @Summary Delete a post by ID
// @Description Deletes a post from the database by its ID
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204 {object} SuccessResponse "No content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
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