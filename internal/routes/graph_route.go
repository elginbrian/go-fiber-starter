package routes

import (
	"bytes"
	"fiber-starter/domain/schema"
	"fiber-starter/internal/di"
	"fiber-starter/pkg/util"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func SetupGraphQLRoute(app *fiber.App, container di.Container) {
	graphqlSchema, err := createGraphQLSchema(container)
	if err != nil {
		log.Printf("Failed to create GraphQL schema: %v", err)
		return
	}

	graphQLHandler := handler.New(&handler.Config{
		Schema: &graphqlSchema,
		Pretty: true,
	})

	app.Post("/api/v1/graphql", graphqlHandler(graphQLHandler))
	app.Get("/graphql/docs", playgroundHandler())
}

func createGraphQLSchema(container di.Container) (graphql.Schema, error) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: util.MergeFields(
			util.ConvertFieldDefinitionMap(schema.NewUserQueryType(container).Fields()),
			util.ConvertFieldDefinitionMap(schema.NewPostQueryType(container).Fields()),
		),
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: util.MergeFields(
			util.ConvertFieldDefinitionMap(schema.NewUserQueryType(container).Fields()),
			util.ConvertFieldDefinitionMap(schema.NewPostMutationType(container).Fields()),
		),
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}

func graphqlHandler(graphQLHandler *handler.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Body()
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
		}

		util.CopyHeaders(c, req)
		rec := util.NewResponseRecorder()
		graphQLHandler.ServeHTTP(rec, req)
		util.CopyResponseHeaders(c, rec)

		return c.Status(rec.StatusCode).Send(rec.Body.Bytes())
	}
}

func playgroundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create request")
		}

		util.CopyHeaders(c, req)
		rec := util.NewResponseRecorder()
		playground.Handler("GraphQL Playground", "/api/v1/graphql").ServeHTTP(rec, req)
		util.CopyResponseHeaders(c, rec)

		return c.Status(rec.StatusCode).Send(rec.Body.Bytes())
	}
}