package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/NickVasky/MaTrOS/pkg/job"
	"github.com/segmentio/kafka-go"
)

type RunnerService struct {
	cfg   *config.RunnerServiceConfig
	kafka *kafka.Reader
}

func NewRunnerService(cfg *config.RunnerServiceConfig, kfk *kafka.Reader) *RunnerService {
	service := new(RunnerService)
	service.cfg = cfg
	service.kafka = kfk
	return service
}

func (s *RunnerService) Serve() {
	wg := &sync.WaitGroup{}
	msgs := make(chan job.Job, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for kafka
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, msgs chan<- job.Job) {
		defer wg.Done()
		for {
			msg, err := s.kafka.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Kafka shutdown")
					return
				}
				log.Println("Kafka read err: ", err)
				continue
			}
			job := job.Job{}
			err = json.Unmarshal(msg.Value, &job)
			if err != nil {
				log.Println("Kafka msg unmarshalling err: ", err)
				continue
			}
			msgs <- job
		}
	}(ctx, wg, msgs)

	// process job
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup, msgs <-chan job.Job) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Queue reader shutdown")
				return
			case msg := <-msgs:
				runJob(ctx, msg)
			}
		}

	}(ctx, wg, msgs)

	wg.Wait()
}

func runJob(ctx context.Context, j job.Job) {
	// it should run jobs via api (http.Client)
	// TBD
}
