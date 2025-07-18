package service

import (
	"fmt"
	"os"
	"time"

	"github.com/emersion/go-imap/v2"
	"gopkg.in/yaml.v3"
)

type TriggersConfigYaml struct {
	Triggers map[string]TriggerYaml `yaml:"triggers"`
}

type TriggerYaml struct {
	BotID     string            `yaml:"bot_id" json:"bot_id"`
	ProcessID string            `yaml:"process_id" json:"process_id"`
	Headers   []HeaderFieldYaml `yaml:"headers" json:"headers"`
	Subject   []string          `yaml:"subject" json:"subject"`
	Body      []string          `yaml:"body" json:"body"`
}

type HeaderFieldYaml struct {
	Key   string `yaml:"key" json:"key"`
	Value string `yaml:"value" json:"value"`
}

func LoadTriggers(path string) *TriggersConfigYaml {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	triggers, err := parseTriggersYaml(data)
	if err != nil {
		panic(err)
	}
	return triggers
}

func parseTriggersYaml(yamlData []byte) (*TriggersConfigYaml, error) {
	triggers := new(TriggersConfigYaml)
	err := yaml.Unmarshal(yamlData, triggers)
	if err != nil {
		return triggers, err
	}
	return triggers, nil
}

func (t *TriggerYaml) Validate() error {
	if t.BotID == "" {
		return fmt.Errorf("no `bot_id` field")
	}
	if t.ProcessID == "" {
		return fmt.Errorf("no `process_id` field")
	}
	if len(t.Headers) == 0 && len(t.Subject) == 0 && len(t.Body) == 0 {
		return fmt.Errorf("no search criteria")
	}
	return nil
}

func (t *TriggerYaml) BuildSearchCriteria() *imap.SearchCriteria {
	searchCriteria := new(imap.SearchCriteria)

	headers := make([]imap.SearchCriteriaHeaderField, 0)
	for _, v := range t.Headers {
		headers = append(headers, imap.SearchCriteriaHeaderField{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	searchCriteria.Body = t.Body
	searchCriteria.Header = headers
	searchCriteria.Since = time.Now().Add(-48 * time.Hour)

	return searchCriteria
}
