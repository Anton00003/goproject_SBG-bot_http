package main

import (
	"goproject_SBG-bot_http/api"
	"goproject_SBG-bot_http/api_http"
	"goproject_SBG-bot_http/data"
	"goproject_SBG-bot_http/repository"
	"goproject_SBG-bot_http/service"
)

func main() {

	databasereader := data.NewDatabaseReader()
	repo := repository.New(databasereader)
	srv := service.New(repo)
	api_http := api_http.New(srv)
	api := api.New(srv)
	go api.Run(srv)
	api_http.Run()

}
