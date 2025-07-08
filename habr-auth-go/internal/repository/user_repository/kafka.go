package user_repository

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

func (r *Repository) SendVerificationCode(email string, code int) error {
	data := map[string]int{email: code}
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return err
	}
	_, _, err = r.kafkaExplorer.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: r.kafkaExplorer.Topic[0],
		Value: sarama.ByteEncoder(msg)})
	if err != nil {
		log.Printf("Failed to SendVerificationCode: %v", err)
		return err
	}
	return nil
}
