package schema

import (
	di "fiber-starter/internal/di"

	"github.com/graphql-go/graphql"
)

var PostType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.ID},
		"userId":    &graphql.Field{Type: graphql.ID},
		"caption":   &graphql.Field{Type: graphql.String},
		"imageURL":  &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.String},
		"updatedAt": &graphql.Field{Type: graphql.String},
	},
})

func NewPostQueryType(container di.Container) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getAllPosts": &graphql.Field{
				Type:    graphql.NewList(PostType),
				Resolve: container.PostResolver.GetAllPosts,
			},
			"getPostByID": &graphql.Field{
				Type: PostType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.ID},
				},
				Resolve: container.PostResolver.GetPostByID,
			},
			"getPostsByUserID": &graphql.Field{
				Type: graphql.NewList(PostType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.ID},
				},
				Resolve: container.PostResolver.GetPostsByUserID,
			},
		},
	})
}

func NewPostMutationType(container di.Container) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPost": &graphql.Field{
				Type: PostType,
				Args: graphql.FieldConfigArgument{
					"userId":  &graphql.ArgumentConfig{Type: graphql.ID},
					"caption": &graphql.ArgumentConfig{Type: graphql.String},
					//"image":   &graphql.ArgumentConfig{Type: graphql.Upload}, 
				},
				Resolve: container.PostResolver.CreatePost,
			},
			"updatePostCaption": &graphql.Field{
				Type: PostType,
				Args: graphql.FieldConfigArgument{
					"id":      &graphql.ArgumentConfig{Type: graphql.ID},
					"caption": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: container.PostResolver.UpdatePostCaption,
			},
			"deletePost": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.ID},
				},
				Resolve: container.PostResolver.DeletePost,
			},
		},
	})
}