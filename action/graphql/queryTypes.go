package graphql

import (
	"github.com/graphql-go/graphql"
	graphqlType "hcc/piano/action/graphql/type"
	"hcc/piano/driver"
	"hcc/piano/lib/logger"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"telegraf": &graphql.Field{
				Type:        graphqlType.TelegrafType,
				Description: "Get all cpu usage data",
				Args: graphql.FieldConfigArgument{
					"metric": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"subMetric": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"period": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"aggregateType": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"duration": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: cpu")
					return driver.GetInfluxData(params.Args)
				},
			},
		},
	})
