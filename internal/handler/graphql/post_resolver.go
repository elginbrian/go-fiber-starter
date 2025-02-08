package handler

import (
	"fmt"
	"log"

	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fiber-starter/pkg/response"
	"fiber-starter/pkg/util"

	"github.com/graphql-go/graphql"
)

type PostResolver struct {
	postService contract.IPostService
	authService contract.IAuthService
}

func NewPostResolver(postService contract.IPostService, authService contract.IAuthService) *PostResolver {
	return &PostResolver{postService, authService}
}

func (r *PostResolver) GetAllPosts(p graphql.ResolveParams) (interface{}, error) {
	posts, err := r.postService.FetchAllPosts()
	if err != nil {
		log.Println("Error fetching posts:", err)
		return nil, fmt.Errorf("failed to fetch posts")
	}

	if posts == nil {
		posts = []entity.Post{} 
	}

	log.Println("Fetched posts:", posts)

	postResponses := make([]response.Post, 0, len(posts))
	for _, post := range posts {
		postResponses = append(postResponses, util.MapToSinglePostResponse(post))
	}

	return postResponses, nil
}

func (r *PostResolver) GetPostByID(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(string)

	post, err := r.postService.FetchPostByID(id)
	if err != nil {
		log.Println("Post not found:", err)
		return nil, fmt.Errorf("post not found")
	}

	return util.MapToSinglePostResponse(post), nil
}

func (r *PostResolver) GetPostsByUserID(p graphql.ResolveParams) (interface{}, error) {
	userID, _ := p.Args["userId"].(string)

	posts, err := r.postService.FetchPostsByUserID(userID)
	if err != nil {
		log.Println("Error fetching posts for user:", err)
		return nil, fmt.Errorf("no posts found for this user")
	}

	if posts == nil {
		posts = []entity.Post{} 
	}

	postResponses := make([]response.Post, 0, len(posts))
	for _, post := range posts {
		postResponses = append(postResponses, util.MapToSinglePostResponse(post))
	}
	return postResponses, nil
}

func (r *PostResolver) CreatePost(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	token, err := util.ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := r.authService.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	caption, _ := p.Args["caption"].(string)
	// imageFile, _ := p.Args["image"].(multipart.File)

	if caption == "" {
		return nil, fmt.Errorf("caption cannot be empty")
	}

	// var imageURL string
	// if imageFile != nil {
	// 	imageURL, err = util.UploadPostImage(imageFile, user.ID, "./uploads/posts/")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to upload image")
	// 	}
	// }

	post := entity.Post{
		UserID:   user.ID,
		Caption:  caption,
		// ImageURL: imageURL,
	}

	createdPost, err := r.postService.CreatePost(post)
	if err != nil {
		return nil, fmt.Errorf("error creating post")
	}

	return util.MapToSinglePostResponse(createdPost), nil
}

func (r *PostResolver) UpdatePostCaption(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	token, err := util.ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := r.authService.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	postID, _ := p.Args["id"].(string)
	caption, _ := p.Args["caption"].(string)

	if caption == "" {
		return nil, fmt.Errorf("caption cannot be empty")
	}

	post, err := r.postService.FetchPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	if post.UserID != user.ID {
		return nil, fmt.Errorf("unauthorized to update this post")
	}

	post.Caption = caption
	updatedPost, err := r.postService.UpdatePost(postID, post)
	if err != nil {
		return nil, fmt.Errorf("error updating post")
	}

	return util.MapToSinglePostResponse(updatedPost), nil
}

func (r *PostResolver) DeletePost(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	token, err := util.ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := r.authService.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	postID, _ := p.Args["id"].(string)

	post, err := r.postService.FetchPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	if post.UserID != user.ID {
		return nil, fmt.Errorf("unauthorized to delete this post")
	}

	err = r.postService.DeletePost(postID)
	if err != nil {
		return nil, fmt.Errorf("error deleting post")
	}

	return true, nil
}