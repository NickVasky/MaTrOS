package main

import (
	"context"
	"fmt"

	"github.com/NickVasky/MaTrOS/runner/client"
	"github.com/NickVasky/MaTrOS/runner/service"
	"github.com/NickVasky/MaTrOS/shared/config"
	"github.com/NickVasky/MaTrOS/shared/queue"
)

func main() {
	cfg := config.NewRunnerServiceConfig()

	kfk := queue.NewConsumer(cfg.Kafka)
	api, err := client.NewBotApiClient(cfg.Orch)
	if err != nil {
		panic(fmt.Errorf("Unable to init Orch client: %v", err))
	}

	s := service.NewRunnerService(cfg, kfk, api)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Serve(ctx)
}
