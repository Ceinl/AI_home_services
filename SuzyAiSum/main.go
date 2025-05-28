package main

import (
	"github.com/Cein13/SuzyAiSum/env"
	"github.com/Cein13/SuzyAiSum/server"
)

func main() {
	cfg := env.LoadConfig() 

	srv := server.NewServer(cfg.API_KEY,cfg.PORT)
	srv.RegisterApi()
	srv.StartWithGracefulShutdown()

}

