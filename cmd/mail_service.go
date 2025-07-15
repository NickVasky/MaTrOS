package main

import (
	"context"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/NickVasky/MaTrOS/mailclient"
	"github.com/NickVasky/MaTrOS/service"
)

func main() {
	cfg := config.NewConfig()

	client, err := mailclient.ConnectToIMAP(&cfg.Mail)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	triggers := service.LoadTriggers("triggers.yaml")

	service, err := service.NewMailListernerService(client, cfg, triggers)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctxCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	service.ListenForMail(ctxCancel)
}
