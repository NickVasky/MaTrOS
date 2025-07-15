package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/NickVasky/MaTrOS/mailclient"
	"github.com/emersion/go-imap/v2"
)

type ListenerJob struct {
	TriggerId string
	Trigger   TriggerYaml
	MailUID   imap.UID
}

type mailListenerService struct {
	cfg      *config.ServiceConfig
	triggers *TriggersConfigYaml
	client   *mailclient.MailClient
}

type TriggerCriteria struct {
	Headers []imap.SearchCriteriaHeaderField
}

func NewMailListernerService(client *mailclient.MailClient, cfg *config.ServiceConfig, triggers *TriggersConfigYaml) (*mailListenerService, error) {
	service := new(mailListenerService)
	service.cfg = cfg
	service.triggers = triggers

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
	go func(ctx context.Context, mail <-chan ListenerJob) {
		for {
			select {
			case m := <-mail:
				fmt.Printf("Got message:\n%v\n", m)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, mailCh)

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
