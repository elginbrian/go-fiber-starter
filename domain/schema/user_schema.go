package schema

import (
	di "fiber-starter/internal/di"

	"github.com/graphql-go/graphql"
)

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.ID},
		"username":  &graphql.Field{Type: graphql.String},
		"email":     &graphql.Field{Type: graphql.String},
		"bio":       &graphql.Field{Type: graphql.String},
		"imageURL":  &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.String},
		"updatedAt": &graphql.Field{Type: graphql.String},
	},
})

func NewUserQueryType(container di.Container) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getAllUsers": &graphql.Field{
				Type:    graphql.NewList(UserType),
				Resolve: container.UserResolver.GetAllUsers,
			},
			"getUserByID": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.ID},
				},
				Resolve: container.UserResolver.GetUserByID,
			},
			"searchUsers": &graphql.Field{
				Type: graphql.NewList(UserType),
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: container.UserResolver.SearchUsers,
			},
		},
	})
}

func NewUserMutationType(container di.Container) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"updateUser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"id":       &graphql.ArgumentConfig{Type: graphql.ID},
					"username": &graphql.ArgumentConfig{Type: graphql.String},
					"bio":      &graphql.ArgumentConfig{Type: graphql.String},
					// "image":    &graphql.ArgumentConfig{Type: graphql.Upload},
				},
				Resolve: container.UserResolver.UpdateUser,
			},
		},
	})
}