package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/NickVasky/MaTrOS/internal/maillistenerservice/mailclient"
	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/emersion/go-imap/v2"
	"github.com/segmentio/kafka-go"
)

type Cacher interface {
	Get(key imap.UID) (interface{}, bool)
	Set(key imap.UID, value interface{})
	Delete(key imap.UID)
	Has(key imap.UID) bool
}

type ListenerJob struct {
	TriggerId string      `json:"trigger_id"`
	Trigger   TriggerYaml `json:"trigger"`
	MailUID   imap.UID    `json:"mail_uid"`
}

type mailListenerService struct {
	cfg      *config.ServiceConfig
	triggers *TriggersConfigYaml
	client   *mailclient.MailClient
	kafka    *kafka.Writer
	cache    Cacher
}

type TriggerCriteria struct {
	Headers []imap.SearchCriteriaHeaderField
}

func NewMailListernerService(client *mailclient.MailClient, cfg *config.ServiceConfig, triggers *TriggersConfigYaml, kfk *kafka.Writer, cache Cacher) (*mailListenerService, error) {
	service := new(mailListenerService)
	service.cfg = cfg
	service.triggers = triggers
	service.kafka = kfk
	service.cache = cache

	service.client = client

	return service, nil
}

func (s *mailListenerService) ListenForMail(ctx context.Context) {
	wg := &sync.WaitGroup{}
	var bufferSize int = 1
	if len(s.triggers.Triggers) > 1 {
		bufferSize = len(s.triggers.Triggers)
	}
	mailCh := make(chan ListenerJob, bufferSize)
	defer close(mailCh)

	// Results gathering goroutine
	s.listenResults(ctx, mailCh)

	// Workers
	for triggerId, trigger := range s.triggers.Triggers {
		job := ListenerJob{
			TriggerId: triggerId,
			Trigger:   trigger,
		}
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, job ListenerJob) {
			log.Println("Listener created for trigger: ", job)
			ticker := time.NewTicker(s.cfg.Mail.PollingInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					s.searchMail(job, mailCh)
				case <-ctx.Done():
					wg.Done()
					return
				}
			}

		}(ctx, wg, job)
	}

	wg.Wait()
}

func (s *mailListenerService) searchMail(job ListenerJob, result chan<- ListenerJob) {
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

func (s *mailListenerService) listenResults(ctx context.Context, result <-chan ListenerJob) {
	ticker := time.NewTicker(50 * time.Millisecond)
	msgBuffer := make([]kafka.Message, 0, 10)
	go func(ctx context.Context, mail <-chan ListenerJob) {
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
				if s.cache.Has(m.MailUID) {
					log.Printf("Cache hit for UID: %v\n", m.MailUID)
					continue
				}
				s.cache.Set(m.MailUID, m)
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

func prepareMessage(job ListenerJob) kafka.Message {
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
