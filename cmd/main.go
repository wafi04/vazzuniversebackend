package main

import (
	"github.com/wafi04/vazzuniversebackend/pkg/server"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
)

func main() {
	log := response.NewLogger()
	server.SetUpAllRoutes()

	log.Info("Service Ready !!!")
}
