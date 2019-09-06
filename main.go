package main

import (
	"GraphQL_piano/pianocheckroot"
	"GraphQL_piano/pianoconfig"
	"GraphQL_piano/pianographql"
	"GraphQL_piano/pianologger"
	"GraphQL_piano/pianomysql"
	"net/http"
)

func main() {
	if !pianocheckroot.CheckRoot() {
		return
	}

	if !pianologger.Prepare() {
		return
	}
	defer pianologger.FpLog.Close()

	err := pianomysql.Prepare()
	if err != nil {
		return
	}
	defer pianomysql.Db.Close()

	http.Handle("/graphql", pianographql.GraphqlHandler)

	pianologger.Logger.Println("Server is running on port " + pianoconfig.HTTPPort)
	err = http.ListenAndServe(":"+pianoconfig.HTTPPort, nil)
	if err != nil {
		pianologger.Logger.Println("Failed to prepare http server!")
	}
}
