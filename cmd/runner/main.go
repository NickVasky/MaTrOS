package main

import (
	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/NickVasky/MaTrOS/pkg/queue"
)

func main() {
	el := config.NewEnvLoader([]string{})
	cfg := config.NewRunnerServiceConfig(el)

	kfk := queue.NewConsumer(cfg.Kafka)
	if kfk != nil {

	}
}
