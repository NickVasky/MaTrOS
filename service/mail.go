package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/joho/godotenv"
)

type MailListenerService struct {
	client          *imapclient.Client
	pollingInterval time.Duration
}

type credentials struct {
	username string
	password string
}

type Cache struct {
	cache map[imap.UID]bool
}

func NewMailListernerService(pollingInterval time.Duration) (*MailListenerService, error) {
	service := new(MailListenerService)

	opts := &imapclient.Options{}
	client, err := imapclient.DialTLS("imap.gmail.com:993", opts)
	if err != nil {
		log.Println("IMAP: Failed to connect")
		return service, err
	}

	creds := getCredentials()

	if err := client.Login(creds.username, creds.password).Wait(); err != nil {
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
	service.pollingInterval = pollingInterval
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
		ticker := time.NewTicker(s.pollingInterval)
		for {
			select {
			case <-ticker.C:
				log.Println("tick!")
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
		Since:  time.Now().Add(-24 * time.Hour),
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

func getCredentials() credentials {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("MAIL_USER")
	password := os.Getenv("MAIL_PASS")

	return credentials{
		username: username,
		password: password,
	}
}
