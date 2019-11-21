package _type

import (
	"github.com/graphql-go/graphql"
)

var TelegrafType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "telegraf",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"series": &graphql.Field{
				Type: graphql.NewList(SeriesType),
			},
		},
	},
)

var SeriesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "seriesType",
		Fields: graphql.Fields{
			"time": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.String,
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
