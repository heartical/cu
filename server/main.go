package main

import (
	"cu/server/api"
	"cu/server/config"
)

func main() {
	api.Init()
	api.StartServer(config.PORT)
}
