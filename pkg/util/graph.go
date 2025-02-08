package util

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type ResponseRecorder struct {
	Headers     http.Header
	Body       bytes.Buffer
	StatusCode int
}

func ExtractTokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return "", fmt.Errorf("missing token")
	}
	return token, nil
}

func Coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func MergeFields(fieldMaps ...graphql.Fields) graphql.Fields {
	merged := graphql.Fields{}
	for _, fields := range fieldMaps {
		for key, value := range fields {
			merged[key] = value
		}
	}
	return merged
}

func ConvertFieldDefinitionMap(fieldMap graphql.FieldDefinitionMap) graphql.Fields {
	fields := graphql.Fields{}
	for key, value := range fieldMap {
		fields[key] = &graphql.Field{
			Type:              value.Type,
			Args:              ConvertArgs(value.Args),
			Resolve:           value.Resolve,
			Description:       value.Description,
			DeprecationReason: value.DeprecationReason,
		}
	}
	return fields
}

func ConvertArgs(args []*graphql.Argument) graphql.FieldConfigArgument {
	fieldConfigArgs := graphql.FieldConfigArgument{}
	for _, arg := range args {
		fieldConfigArgs[arg.Name()] = &graphql.ArgumentConfig{
			Type:         arg.Type,
			DefaultValue: arg.DefaultValue,
			Description:  arg.Description(),
		}
	}
	return fieldConfigArgs
}


func NewResponseRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		Headers:     make(http.Header),
		StatusCode: http.StatusOK,
	}
}

func (r *ResponseRecorder) Header() http.Header {
	return r.Headers
}

func (r *ResponseRecorder) Write(data []byte) (int, error) {
	return r.Body.Write(data)
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

func CopyHeaders(c *fiber.Ctx, req *http.Request) {
	req.Header = make(http.Header)
	c.Request().Header.VisitAll(func(k, v []byte) {
		req.Header.Set(string(k), string(v))
	})
}

func CopyResponseHeaders(c *fiber.Ctx, rec *ResponseRecorder) {
	for k, v := range rec.Headers {
		c.Set(k, v[0])
	}
}