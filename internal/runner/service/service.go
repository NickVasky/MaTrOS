package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/NickVasky/MaTrOS/internal/runner/client"
	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/NickVasky/MaTrOS/pkg/job"
	"github.com/segmentio/kafka-go"
)

type RunnerService struct {
	cfg       *config.RunnerServiceConfig
	kafka     *kafka.Reader
	apiclient *client.BotApiClient
}

func NewRunnerService(cfg *config.RunnerServiceConfig, kfk *kafka.Reader, api *client.BotApiClient) *RunnerService {
	service := new(RunnerService)
	service.cfg = cfg
	service.kafka = kfk
	service.apiclient = api
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
				err := s.runJob(ctx, msg)
				if err != nil {
					log.Println("WARNING! Failed job:", msg)
					// TODO - Alerts and retry logic goes somewhere here
				}
			}
		}

	}(ctx, wg, msgs)

	wg.Wait()
}

func (s *RunnerService) runJob(_ context.Context, j job.Job) error {
	projects, err := s.apiclient.GetProjects()
	if err != nil {
		log.Println("runJob - GetProjects() error: ", err)
		return err
	}
	foundProject, err := projects.GetProjectByName(j.Trigger.ProcessID)
	if err != nil {
		log.Println("runJob - Project not found. ", err)
		return err
	}
	robots, err := s.apiclient.GetRobots()
	if err != nil {
		log.Println("runJob - GetRobots() error: ", err)
		return err
	}
	foundRobot, err := robots.GetRobotByName(j.Trigger.BotID)
	if err != nil {
		log.Println("runJob - Robot not found. ", err)
		return err
	}

	err = s.apiclient.PutRobotStartAsync(foundRobot.ID, foundProject.ID)
	if err != nil {
		log.Println("runJob - PutRobotStartAsync() error: ", err)
		return err
	}

	return nil
}
