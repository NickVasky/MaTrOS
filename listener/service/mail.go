package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NickVasky/MaTrOS/listener/cache"
	"github.com/NickVasky/MaTrOS/listener/mailclient"
	"github.com/NickVasky/MaTrOS/shared/config"
	"github.com/NickVasky/MaTrOS/shared/job"
	"github.com/emersion/go-imap/v2"
	"github.com/segmentio/kafka-go"
)

type mailListenerService struct {
	cfg      *config.MailListenerServiceConfig
	triggers *job.TriggersConfigYaml
	client   *mailclient.MailClient
	kafka    *kafka.Writer
	cache    cache.Cacher
}

func NewMailListernerService(client *mailclient.MailClient, cfg *config.MailListenerServiceConfig, triggers *job.TriggersConfigYaml, kfk *kafka.Writer, cache cache.Cacher) (*mailListenerService, error) {
	service := new(mailListenerService)

	service.client = client
	service.cfg = cfg
	service.triggers = triggers
	service.kafka = kfk
	service.cache = cache

	return service, nil
}

func (s *mailListenerService) ListenForMail(ctx context.Context) {
	wg := &sync.WaitGroup{}
	var bufferSize int = 1
	if len(s.triggers.Triggers) > 1 {
		bufferSize = len(s.triggers.Triggers)
	}
	mailCh := make(chan job.Job, bufferSize)
	defer close(mailCh)

	// Results gathering goroutine
	s.listenResults(ctx, mailCh)

	// Workers
	for triggerId, trigger := range s.triggers.Triggers {
		j := job.Job{
			TriggerId: triggerId,
			Trigger:   trigger,
		}
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, j job.Job) {
			log.Println("Listener created for trigger: ", j)
			ticker := time.NewTicker(s.cfg.Mail.PollingInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					s.searchMail(j, mailCh)
				case <-ctx.Done():
					wg.Done()
					return
				}
			}

		}(ctx, wg, j)
	}

	wg.Wait()
}

func (s *mailListenerService) searchMail(job job.Job, result chan<- job.Job) {
	err := s.client.IMAP.Noop().Wait()
	if err != nil {
		log.Println(err)
	}

	searchCriteria := job.Trigger.BuildSearchCriteria()

	searchOpts := &imap.SearchOptions{
		ReturnAll: true,
	}

	msgs, err := s.client.IMAP.UIDSearch(searchCriteria, searchOpts).Wait()
	if err != nil {
		log.Println(err)
	}
	uids := msgs.AllUIDs()

	for _, uid := range uids {
		job.MailUID = uid
		result <- job
	}
}

func (s *mailListenerService) listenResults(ctx context.Context, result <-chan job.Job) {
	ticker := time.NewTicker(50 * time.Millisecond)
	msgBuffer := make([]kafka.Message, 0, 10)
	go func(ctx context.Context, mail <-chan job.Job) {
		for {
			select {
			case <-ticker.C:
				// send anything
				if len(msgBuffer) > 0 {
					log.Println("Timeout for Kafka buffer, sending...")
					err := s.kafka.WriteMessages(ctx, msgBuffer...) // send buffer
					if err != nil {
						log.Println(err)
						continue
					}
					msgBuffer = msgBuffer[:0] // flush buffer
				}

			case m := <-mail:
				key := fmt.Sprintf("mail:%v", m.MailUID)
				cacheCtx, cacheCancel := context.WithTimeout(ctx, 3*time.Second)
				defer cacheCancel()
				hasKey, err := s.cache.Has(cacheCtx, key)
				if err == nil && hasKey {
					log.Printf("Cache hit for UID: %v\n", m.MailUID)
					continue
				}
				if err != nil {
					log.Printf("Cache err: %v\n", err)
				}
				err = s.cache.Set(cacheCtx, key, nil, s.cfg.Redis.TTL)
				if err != nil {
					log.Printf("Cache err: %v\n", err)
				}

				msg := prepareMessage(m)
				msgBuffer = append(msgBuffer, msg)
				log.Printf("Got message:\n%v\n", m)
				if len(msgBuffer) >= 10 {
					log.Println("Kafka buffer is full, sending...")
					err := s.kafka.WriteMessages(ctx, msgBuffer...) // send buffer
					if err != nil {
						log.Println(err)
						continue
					}
					msgBuffer = msgBuffer[:0] // flush buffer
				}

			case <-ctx.Done():
				return
			}
		}
	}(ctx, result)
}

func prepareMessage(job job.Job) kafka.Message {
	var msg kafka.Message

	body, err := json.Marshal(job)
	if err != nil {
		log.Println("Marshalling err: ", err)
	}
	msg.Value = body

	return msg
}

// func fetchLetter() {
// 	var uids []imap.UID
// 	seqset := new(imap.UIDSet)
// 	for _, v := range uids {
// 		seqset.AddNum(v)
// 		fmt.Println(v)
// 	}

// 	bodySection := &imap.FetchItemBodySection{Specifier: imap.PartSpecifierHeader}
// 	fetchOptions := &imap.FetchOptions{
// 		UID:         true,
// 		Flags:       true,
// 		Envelope:    true,
// 		BodySection: []*imap.FetchItemBodySection{bodySection},
// 	}
// 	msgs2, err := s.client.IMAP.Fetch(*seqset, fetchOptions).Collect()
// 	if err != nil {
// 		log.Println("Fetch Err: ", err)
// 		return
// 	}
// 	log.Println("Seq: ", seqset)

// 	log.Println("Printing messages...")
// 	log.Println(len(msgs2))
// 	for _, v := range msgs2 {
// 		fmt.Println(v.Envelope)
// 	}
// }
