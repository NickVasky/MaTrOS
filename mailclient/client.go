package mailclient

import (
	"fmt"
	"log"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type MailClient struct {
	IMAP *imapclient.Client
}

func ConnectToIMAP(cfg *config.MailConfig) (*MailClient, error) {
	mailclient := new(MailClient)

	opts := &imapclient.Options{}
	server_url := fmt.Sprintf("%v:%v", cfg.URL, cfg.Port)
	client, err := imapclient.DialTLS(server_url, opts)
	if err != nil {
		log.Println("IMAP: Failed to connect")
		return mailclient, err
	}

	if err := client.Login(cfg.Username, cfg.Password).Wait(); err != nil {
		log.Println("IMAP: Failed to login")
		return mailclient, err
	}
	log.Println("IMAP: Successfully logged in")

	selectOpts := &imap.SelectOptions{ReadOnly: true}
	_, err = client.Select("INBOX", selectOpts).Wait()
	if err != nil {
		log.Println("IMAP: Selection of mailbox is failed")
		return mailclient, err
	}
	log.Println("IMAP: Folder selected")
	mailclient.IMAP = client

	return mailclient, nil
}

func (m *MailClient) Stop() error {
	err := m.IMAP.Logout().Wait()
	if err != nil {
		return err
	}
	err = m.IMAP.Close()
	if err != nil {
		return err
	}
	return nil
}
