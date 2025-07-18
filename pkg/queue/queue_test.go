package queue

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/NickVasky/MaTrOS/pkg/config"
	"github.com/segmentio/kafka-go"
)

func TestNewProducer(t *testing.T) {
	cfg := &config.KafkaConfig{
		Host:  "localhost:9094",
		Topic: "test_topic",
	}
	ctx := context.Background()
	p := NewProducer(cfg)
	fmt.Println("Connected...")
	msgs := make([]kafka.Message, 0, 10)

	for i := range 10 {
		msgs = append(msgs, kafka.Message{Value: []byte(fmt.Sprintf("Hello numero: %v!", i))})
	}

	err := p.WriteMessages(ctx, msgs...)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	} else {
		fmt.Println("Message sent...")
	}
}

func TestRead(t *testing.T) {
	//kafka.SetLogger(log.New(os.Stdout, "kafka-go: ", log.LstdFlags))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9094"},
		Topic:          "test_topic",
		GroupID:        "my-first-application",
		CommitInterval: 0, // Отключаем автоматический коммит
		StartOffset:    kafka.FirstOffset,
		//Logger:         log.New(os.Stdout, "kafka-go-client: ", log.LstdFlags),
	})

	resetFlag := false
	if resetFlag {
		log.Println("resetting offset")
		msg := kafka.Message{
			Topic:     "test_topic",
			Partition: 0,
			Offset:    0,
		}
		err := reader.CommitMessages(context.Background(), msg)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Reader created")
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n", string(msg.Value))

		// Коммитим оффсет вручную после обработки
		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			panic(err)
		}
	}
}
