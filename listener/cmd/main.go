package main

import (
	"context"

	"github.com/NickVasky/MaTrOS/listener/cache"
	"github.com/NickVasky/MaTrOS/listener/mailclient"
	"github.com/NickVasky/MaTrOS/listener/service"
	"github.com/NickVasky/MaTrOS/shared/config"
	"github.com/NickVasky/MaTrOS/shared/job"
	"github.com/NickVasky/MaTrOS/shared/queue"
)

func main() {
	el := config.NewEnvLoader([]string{})
	cfg := config.NewMailListenerServiceConfig(el)

	kfk := queue.NewProducer(cfg.Kafka)

	cache, err := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.User, cfg.Redis.Password)
	if err != nil {
		panic(err)
	}

	client, err := mailclient.ConnectToIMAP(cfg.Mail)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	triggers := job.LoadTriggers("triggers.yaml")

	service, err := service.NewMailListernerService(client, cfg, triggers, kfk, cache)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctxCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	service.ListenForMail(ctxCancel)
}
