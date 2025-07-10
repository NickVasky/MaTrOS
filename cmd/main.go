package main

import (
	"context"
	"time"

	"github.com/NickVasky/MaTrOS/service"
)

func main() {
	service, err := service.NewMailListernerService(15 * time.Second)
	if err != nil {
		panic(err)
	}

	defer service.Stop()

	ctx := context.Background()
	ctxCancel, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	service.ListenForMail(ctxCancel)
}
