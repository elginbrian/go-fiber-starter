package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var Upload = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Upload",
	Description: "A file upload scalar",
	Serialize:   func(value interface{}) interface{} { return value },
	ParseValue:  func(value interface{}) interface{} { return value },
	ParseLiteral: func(valueAST ast.Value) interface{} {
		return nil
	},
})
