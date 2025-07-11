package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type MailListenerService struct {
	cfg    *config.ServiceConfig
	client *imapclient.Client
}

type TriggerCriteria struct {
	Headers []imap.SearchCriteriaHeaderField
}

func NewMailListernerService(cfg *config.ServiceConfig) (*MailListenerService, error) {
	service := new(MailListenerService)
	service.cfg = cfg

	opts := &imapclient.Options{}
	url := fmt.Sprintf("%v:%v", service.cfg.Mail.URL, service.cfg.Mail.Port)
	client, err := imapclient.DialTLS(url, opts)
	if err != nil {
		log.Println("IMAP: Failed to connect")
		return service, err
	}

	if err := client.Login(service.cfg.Mail.Username, service.cfg.Mail.Password).Wait(); err != nil {
		log.Println("IMAP: Failed to login")
		return service, err
	}
	log.Println("IMAP: Successfully logged in")

	/*
		listCmd := client.List("", "%", &imap.ListOptions{
			ReturnStatus: &imap.StatusOptions{
				NumMessages: true,
				NumUnseen:   true,
			},
		})
		for {
			mbox := listCmd.Next()
			if mbox == nil {
				break
			}
			if mbox.Status != nil {
				log.Printf("Mailbox %q contains %v messages (%v unseen)", mbox.Mailbox, *mbox.Status.NumMessages, *mbox.Status.NumUnseen)
			} else {
				log.Printf("Mailbox %q - Status unavailable", mbox.Mailbox)

			}
		}
		if err := listCmd.Close(); err != nil {
			log.Fatalf("LIST command failed: %v", err)
		}
	*/
	selectOpts := &imap.SelectOptions{ReadOnly: true}
	_, err = client.Select("INBOX", selectOpts).Wait()
	if err != nil {
		log.Println("IMAP: Selection of mailbox is failed")
		return service, err
	}
	log.Println("IMAP: Folder selected")

	service.client = client

	return service, nil
}

func (s *MailListenerService) Stop() {
	s.client.Logout().Wait()
	s.client.Close()
}

func (s *MailListenerService) ListenForMail(ctx context.Context) {
	wg := &sync.WaitGroup{}
	mailCh := make(chan int, 10)
	defer close(mailCh)

	wg.Add(1)
	go func(wg *sync.WaitGroup, mail chan<- int) {
		ticker := time.NewTicker(s.cfg.Mail.PollingInterval)
		for {
			select {
			case <-ticker.C:
				s.fetchMail()
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}(wg, mailCh)

	wg.Wait()
}

func (s *MailListenerService) fetchMail() {
	headers := []imap.SearchCriteriaHeaderField{{
		Key:   "From",
		Value: "mirrayletters@gmail.com"},
	}

	searchCriteria := &imap.SearchCriteria{
		Header: headers,
		Body:   []string{"test", "sub"},
		Since:  time.Now().Add(-48 * time.Hour),
	}

	searchOpts := &imap.SearchOptions{
		ReturnAll: true,
	}

	msgs, err := s.client.UIDSearch(searchCriteria, searchOpts).Wait()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("search data: %v\n", *msgs)
	uids := msgs.AllUIDs()
	fmt.Printf("len: %v\n", len(uids))

	seqset := new(imap.UIDSet)
	for _, v := range uids {
		seqset.AddNum(v)
		fmt.Println(v)
	}

	bodySection := &imap.FetchItemBodySection{Specifier: imap.PartSpecifierHeader}
	fetchOptions := &imap.FetchOptions{
		UID:         true,
		Flags:       true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{bodySection},
	}
	msgs2, err := s.client.Fetch(*seqset, fetchOptions).Collect()
	if err != nil {
		log.Println("Fetch Err: ", err)
		return
	}
	log.Println("Seq: ", seqset)

	log.Println("Printing messages...")
	log.Println(len(msgs2))
	for _, v := range msgs2 {
		fmt.Println(v.Envelope)
	}
}
