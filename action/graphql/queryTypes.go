package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/piano/lib/logger"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{

			"all_subnet": &graphql.Field{
				Type:        graphql.NewList(volumeType),
				Description: "Get all subnet list",
				Args: graphql.FieldConfigArgument{
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_subnet")
					return nil, nil
				},
			},
		},
	})
