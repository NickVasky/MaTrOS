package job

import "github.com/emersion/go-imap/v2"

type Job struct {
	TriggerId string      `json:"trigger_id"`
	Trigger   TriggerYaml `json:"trigger"`
	MailUID   imap.UID    `json:"mail_uid"`
}
