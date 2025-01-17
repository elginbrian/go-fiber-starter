package handler

import (
	"fiber-starter/internal/domain"
	"fiber-starter/internal/service"
	"fiber-starter/pkg/request"
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
// @Success 200 {object} response.GetAllPostsResponse "Successful fetch posts response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
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
// @Param id path string true "Post ID"
// @Success 200 {object} response.GetPostByIDResponse "Successful fetch post response" 
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
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

// GetPostsByUserID godoc
// @Summary Get all posts by a specific user
// @Description Retrieves all posts created by a specific user, including the caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.GetAllPostsResponse "Successful fetch posts by user response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/posts/user/{user_id} [get]
func (h *PostHandler) GetPostsByUserID(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return response.Error(c, "Invalid user ID", fiber.StatusBadRequest)
	}

	posts, err := h.postService.FetchPostsByUserID(userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return response.Error(c, "No posts found for this user", fiber.StatusNotFound)
		}
		return response.Error(c, "Failed to fetch posts", fiber.StatusInternalServerError)
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

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post with an optional image. The caption is required. If an image is provided, it will be uploaded to the server, and the URL will be returned in the response.
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param caption formData string true "Post caption"
// @Param image formData file false "Post image (optional)"
// @Security BearerAuth
// @Success 201 {object} response.CreatePostResponse "Successful image upload response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	if userID == "" {
        return response.Error(c, "Unauthorized", fiber.StatusUnauthorized)
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
// @Summary Update the caption of an existing post
// @Description Updates only the caption of a post. Only the creator of the post is allowed to update it.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body request.UpdatePostRequest true "Request body with updated caption"
// @Security BearerAuth
// @Success 200 {object} response.UpdatePostResponse "Successful update response"
// @Failure 400 {object} response.ErrorResponse "Bad request
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/posts/{id} [put]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)

    id := c.Params("id")
    postID, err := strconv.Atoi(id)
    if err != nil {
        return response.Error(c, "Invalid post ID", fiber.StatusBadRequest)
    }

    existingPost, err := h.postService.FetchPostByID(postID)
    if err != nil {
        return response.Error(c, "Post not found", fiber.StatusNotFound)
    }

    if existingPost.UserID != userID {
        return response.Error(c, "Unauthorized to update this post", fiber.StatusUnauthorized)
    }

    var input request.UpdatePostRequest
    if err := c.BodyParser(&input); err != nil {
        return response.ValidationError(c, "Invalid input")
    }

    if input.Caption == "" {
        return response.ValidationError(c, "Caption cannot be empty")
    }

    existingPost.Caption = input.Caption
    updatedPost, err := h.postService.UpdatePost(postID, existingPost)
    if err != nil {
        return response.Error(c, "Failed to update post", fiber.StatusInternalServerError)
    }

    postResponse := domain.PostResponse{
        ID:        updatedPost.ID,
        UserID:    updatedPost.UserID,
        Caption:   updatedPost.Caption,
        ImageURL:  updatedPost.ImageURL,
        CreatedAt: updatedPost.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: updatedPost.UpdatedAt.Format("2006-01-02 15:04:05"),
    }

    return response.Success(c, postResponse, fiber.StatusOK)
}


// DeletePost godoc
// @Summary Delete a post
// @Description Deletes a post by its ID. Only the creator of the post is authorized to delete it.
// @Tags posts
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 204 {object} response.DeletePostResponse "Successful delete post response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

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

// SearchPosts godoc
// @Summary Search posts by title or content
// @Description Retrieves posts that match the search query.
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} response.SearchPostsResponse "Successful search response"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameter"
// @Router /api/search/posts [get]
func (h *PostHandler) SearchPosts(c *fiber.Ctx) error {
	query := c.Query("query")

	if query == "" {
		return response.Error(c, "Query parameter is required", fiber.StatusBadRequest)
	}

	posts, err := h.postService.SearchPosts(query)
	if err != nil {
		return response.Error(c, "No posts found", fiber.StatusNotFound)
	}

	var postResponses []domain.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, domain.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Caption:   post.Caption,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.Success(c, postResponses)
}