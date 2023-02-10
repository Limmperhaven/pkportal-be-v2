package main

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/middlewares"
	"github.com/Limmperhaven/pkportal-be-v2/internal/domain"
	"github.com/Limmperhaven/pkportal-be-v2/internal/server"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	err = stpg.InitConnect(&cfg.Postgres)
	if err != nil {
		log.Fatalf("error initializing database: %s", err.Error())
	}
	uc := domain.NewUsecase()
	c := controllers.NewController(uc)
	m := middlewares.NewMiddlewareStorage()
	srv := server.NewServer(&cfg.Server, c, m)
	srv.Run()
}
