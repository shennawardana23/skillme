package main

import (
	"fixture/pkg/config"
	"fixture/pkg/server"
)

func main() {
	cfg := config.Load()
	server.Start(cfg)
}
