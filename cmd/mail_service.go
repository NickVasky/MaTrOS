package main

import (
	"context"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/NickVasky/MaTrOS/service"
)

func main() {
	cfg := config.NewConfig()
	service, err := service.NewMailListernerService(cfg)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	ctx := context.Background()
	ctxCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	service.ListenForMail(ctxCancel)
}
