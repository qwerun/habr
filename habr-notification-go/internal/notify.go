package internal

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/qwerun/habr-notification-go/pkg/kafka"
	"log"
	"time"
)

type Consumer struct {
	*kafka.Consumer
}

type Handler struct {
}

func NewConsumer(consumer *kafka.Consumer) *Consumer {
	return &Consumer{consumer}
}

func (c *Consumer) Notify(ctx context.Context) error {
	handler := &Handler{}

	for {
		select {
		case <-ctx.Done():
			log.Println("kafka ctx Done")
			return nil
		default:
		}

		err := c.Group.Consume(ctx, c.Topics, handler)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Shutdown by context")
				return nil
			}

			log.Printf("kafka consume error: %v\n", err)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
