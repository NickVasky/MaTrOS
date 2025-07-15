package main

import (
	"context"

	"github.com/NickVasky/MaTrOS/cache"
	"github.com/NickVasky/MaTrOS/config"
	"github.com/NickVasky/MaTrOS/mailclient"
	"github.com/NickVasky/MaTrOS/queue"
	"github.com/NickVasky/MaTrOS/service"
)

func main() {
	cfg := config.NewConfig()

	kfk := queue.NewProducer(&cfg.Kafka)
	cache := cache.NewInMemoryCache()

	client, err := mailclient.ConnectToIMAP(&cfg.Mail)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	triggers := service.LoadTriggers("triggers.yaml")

	service, err := service.NewMailListernerService(client, cfg, triggers, kfk, cache)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctxCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	service.ListenForMail(ctxCancel)
}
