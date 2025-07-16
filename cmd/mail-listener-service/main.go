package main

import (
	"context"

	"github.com/NickVasky/MaTrOS/internal/maillistenerservice/cache"
	"github.com/NickVasky/MaTrOS/internal/maillistenerservice/mailclient"
	"github.com/NickVasky/MaTrOS/internal/maillistenerservice/service"
	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/NickVasky/MaTrOS/pkg/queue"
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
