package kafka

import (
	"github.com/IBM/sarama"
	"github.com/qwerun/habr-notification-go/internal/notify"
	"log"
)

type ConsumerHandler struct {
	Handler notify.NotificationHandler
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer group setup")
	return nil
}
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer group cleanup")
	return nil
}
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.Handler.HandleNotification(message.Value); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
