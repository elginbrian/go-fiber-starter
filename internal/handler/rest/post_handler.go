package handler

import (
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/request"
	"fiber-starter/pkg/response"
	"fiber-starter/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postService contract.IPostService
	authService contract.IAuthService
}

func NewPostHandler(postService contract.IPostService, authService contract.IAuthService) *PostHandler {
	return &PostHandler{
        postService: postService,
        authService: authService,
    }
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Get a list of all posts, along with details like the user who created them, the caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Success 200 {object} response.GetAllPostsResponse "Successful fetch posts response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	posts, err := h.postService.FetchAllPosts()
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}
	return response.Success(c, util.MapToPostResponse(posts), fiber.StatusOK)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Get a post by its unique ID, including the caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.GetPostByIDResponse "Successful fetch post response" 
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.postService.FetchPostByID(id)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	postResponse := util.MapToPostResponse([]entity.Post{post})[0]
	return response.Success(c, postResponse, fiber.StatusOK)
}

// GetPostsByUserID godoc
// @Summary Get all posts by a specific user
// @Description Get all posts made by a specific user, including the caption, image URL, and timestamps.
// @Tags posts
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} response.GetAllPostsResponse "Successful fetch posts by user response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/user/{user_id} [get]
func (h *PostHandler) GetPostsByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	posts, err := h.postService.FetchPostsByUserID(userID)
	if err != nil {
		return response.Error(c, "Failed to fetch posts", fiber.StatusInternalServerError)
	}
	return response.Success(c, util.MapToPostResponse(posts), fiber.StatusOK)
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with a caption. Optionally, you can upload an image. If an image is uploaded, its URL will be returned in the response. Requires JWT authentication.
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param caption formData string true "Post caption"
// @Param image formData file false "Post image (optional)"
// @Security BearerAuth
// @Success 201 {object} response.CreatePostResponse "Successful image upload response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	caption := c.FormValue("caption")
	if caption == "" {
		return response.ValidationError(c, "Caption is required")
	}

	imageURL, err := util.UploadPostImage(c, user.ID, "./uploads/posts/")
	if err != nil {
		return response.Error(c, "Failed to upload image", fiber.StatusInternalServerError)
	}

	post := entity.Post{
		UserID:   user.ID,
		Caption:  caption,
		ImageURL: imageURL, 
	}

	createdPost, err := h.postService.CreatePost(post)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	postResponse := util.MapToPostResponse([]entity.Post{createdPost})[0]
	return response.Success(c, postResponse, fiber.StatusCreated)
}

// UpdatePost godoc
// @Summary Update an existing post's caption
// @Description Update only the caption of an existing post. Only the post creator is allowed to make this change. Requires JWT authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body request.UpdatePostRequest true "Request body with updated caption"
// @Security BearerAuth
// @Success 200 {object} response.UpdatePostResponse "Successful update response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{id} [patch]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	id := c.Params("id")
	existingPost, err := h.postService.FetchPostByID(id)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	if existingPost.UserID != user.ID {
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
	updatedPost, err := h.postService.UpdatePost(id, existingPost)
	if err != nil {
		return response.Error(c, "Failed to update post", fiber.StatusInternalServerError)
	}

	postResponse := util.MapToPostResponse([]entity.Post{updatedPost})[0]
	return response.Success(c, postResponse, fiber.StatusOK)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post by its ID. Only the post creator is allowed to delete it. Requires JWT authentication.
// @Tags posts
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 204 {object} response.DeletePostResponse "Successful delete post response"
// @Failure 400 {object} response.ErrorResponse "Bad request"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	user, err := util.GetUserFromToken(c, h.authService)
	if err != nil {
		return err
	}

	id := c.Params("id")
	post, err := h.postService.FetchPostByID(id)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}

	if post.UserID != user.ID {
		return response.Error(c, "Unauthorized to delete this post", fiber.StatusUnauthorized)
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		return response.Error(c, "Post not found", fiber.StatusNotFound)
	}
	return response.Success(c, fiber.Map{"message": "Post deleted successfully"}, fiber.StatusNoContent)
}

// SearchPosts godoc
// @Summary Search posts
// @Description Search for posts that match a given query, such as a keyword in the caption or content.
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {array} response.SearchPostsResponse "Successful search response"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameter"
// @Router /search/posts [get]
func (h *PostHandler) SearchPosts(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return response.Error(c, "Query parameter is required", fiber.StatusBadRequest)
	}

	posts, err := h.postService.SearchPosts(query)
	if err != nil {
		return response.Error(c, "No posts found", fiber.StatusNotFound)
	}

	return response.Success(c, util.MapToPostResponse(posts))
}