package main

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/client"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/middlewares"
	"github.com/Limmperhaven/pkportal-be-v2/internal/domain"
	"github.com/Limmperhaven/pkportal-be-v2/internal/server"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
	"github.com/Limmperhaven/pkportal-be-v2/internal/worker"
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
	mail := client.NewMailClient(&cfg.SMTP)
	s3, err := client.InitMinio(&cfg.S3)
	if err != nil {
		log.Fatalf("error initializing minio client: %s", err.Error())
	}
	uc := domain.NewUsecase(mail, s3)
	c := controllers.NewController(uc)
	m := middlewares.NewMiddlewareStorage()
	srv := server.NewServer(&cfg.Server, c, m)
	workerPool := worker.NewPool(mail, s3)
	workerPool.AddWorker(worker.NotificationWorker)
	workerPool.Start()
	srv.Run()
}
