package queue

import (
	"time"

	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/segmentio/kafka-go"
)

func NewProducer(cfg *config.KafkaConfig) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP([]string{cfg.Host}...),
		Topic:        cfg.Topic,
		RequiredAcks: kafka.RequireAll,
		BatchSize:    10,
		BatchTimeout: 10 * time.Millisecond,
	}
	return w
}

func NewConsumer(cfg *config.KafkaConfig) *kafka.Reader {
	rcfg := kafka.ReaderConfig{
		Brokers: []string{cfg.Host},
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	}
	r := kafka.NewReader(rcfg)
	return r
}
