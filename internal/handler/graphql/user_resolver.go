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

type UserResolver struct {
	userService contract.IUserService
	authService contract.IAuthService
}

func NewUserResolver(userService contract.IUserService, authService contract.IAuthService) *UserResolver {
	return &UserResolver{userService, authService}
}

func (r *UserResolver) GetAllUsers(p graphql.ResolveParams) (interface{}, error) {
	users, err := r.userService.FetchAllUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, fmt.Errorf("failed to fetch users")
	}

	if users == nil {
		users = []entity.User{}
	}
	
	userResponses := make([]response.User, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, util.MapToUserResponse(user))
	}

	return userResponses, nil
}

func (r *UserResolver) GetUserByID(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(string)

	user, err := r.userService.FetchUserByID(id)
	if err != nil {
		log.Println("User not found:", err)
		return nil, fmt.Errorf("user not found")
	}

	return util.MapToUserResponse(user), nil
}

func (r *UserResolver) SearchUsers(p graphql.ResolveParams) (interface{}, error) {
	query, _ := p.Args["query"].(string)
	if query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	users, err := r.userService.SearchUsers(query)
	if err != nil {
		if err.Error() == "no users found" {
			return []response.User{}, nil
		}
		log.Println("Error searching users:", err)
		return nil, fmt.Errorf("failed to search users")
	}

	userResponses := make([]response.User, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, util.MapToUserResponse(user))
	}

	return userResponses, nil
}

func (r *UserResolver) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	token, err := util.ExtractTokenFromContext(ctx)
	if err != nil {
		log.Println("Unauthorized:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := r.authService.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	username, _ := p.Args["username"].(string)
	bio, _ := p.Args["bio"].(string)
	// imageFile, _ := p.Args["image"].(multipart.File)

	var imageURL string
	// if imageFile != nil {
	// 	imageURL, err = util.UploadProfileImage(imageFile, user.ID, "./uploads/profiles/")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to upload image: %v", err)
	// 	}
	// }

	if username != "" && (len(username) < 3 || len(username) > 50) {
		return nil, fmt.Errorf("username must be between 3 and 50 characters")
	}

	updatedUser := entity.User{
		ID:        user.ID,
		Name:      util.Coalesce(username, user.Name),
		Email:     user.Email,
		Bio:       util.Coalesce(bio, user.Bio),
		ImageURL:  util.Coalesce(imageURL, user.ImageURL),
		CreatedAt: user.CreatedAt,
	}

	updatedUser, err = r.userService.UpdateUser(user.ID, updatedUser)
	if err != nil {
		log.Println("Error updating user:", err)
		return nil, fmt.Errorf("error updating user")
	}

	return util.MapToUserResponse(updatedUser), nil
}