package internal

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/qwerun/habr-notification-go/pkg/kafka"
	"log"
	"time"
)

type Consumer struct {
	*kafka.KafkaConsumerExplorer
}

type Handler struct {
}

func NewConsumer(consumer *kafka.KafkaConsumerExplorer) *Consumer {
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

		err := c.ConsumerGroup.Consume(ctx, c.Topics, handler)
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
	for message := range claim.Messages() {
		// отправка на почту была бы здесь в отдельной функции
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d, offset = %d",
			string(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset)

		session.MarkMessage(message, "")
	}
	return nil
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer group setup")
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer group cleanup")
	return nil
}
