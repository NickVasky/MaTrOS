package queue

import (
	"fmt"

	"github.com/NickVasky/MaTrOS/config"
	"github.com/segmentio/kafka-go"
)

func NewProducer(cfg *config.KafkaConfig) *kafka.Writer {
	addr := fmt.Sprintf("%v:%v", cfg.URL, cfg.Port)
	w := &kafka.Writer{
		Addr:         kafka.TCP([]string{addr}...),
		Topic:        cfg.Topic,
		RequiredAcks: kafka.RequireAll,
	}
	return w
}
