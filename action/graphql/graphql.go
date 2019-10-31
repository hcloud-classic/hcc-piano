package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

//Schema - cgs
var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryTypes,
	},
)

//GraphqlHandler - cgs
var GraphqlHandler = handler.New(&handler.Config{
	Schema:   &Schema,
	Pretty:   true,
	GraphiQL: true,
})
