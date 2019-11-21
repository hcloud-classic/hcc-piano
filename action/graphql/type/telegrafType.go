package graphqlType

import (
	"github.com/graphql-go/graphql"
)

// TelegrafType : GraphQL type of Telegraf
var TelegrafType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "telegraf",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"data": &graphql.Field{
				Type: graphql.NewList(SeriesType),
			},
		},
	},
)

// SeriesType : GraphQL type of Series
var SeriesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "seriesType",
		Fields: graphql.Fields{
			"x": &graphql.Field{
				Type: graphql.Int,
			},
			"y": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

//var CpuType = graphql.NewObject(
//	graphql.ObjectConfig{
//		Name: "Telegraf",
//		Fields: graphql.Fields{
//			"metric": &graphql.Field{
//				Type: graphql.String,
//			},
//			"subMetric": &graphql.Field{
//				Type: graphql.String,
//			},
//			"period": &graphql.Field{
//				Type: graphql.Int,
//			},
//			"aggregateType": &graphql.Field{
//				Type: graphql.String,
//			},
//			"duration": &graphql.Field{
//				Type: graphql.String,
//			},
//			"uuid": &graphql.Field{
//				Type: graphql.String,
//			},
//		},
//	},
//)
