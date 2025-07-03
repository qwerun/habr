package notify

import "log"

type EmailHandler struct {
	// smpt
}

func (h *EmailHandler) HandleNotification(msg []byte) error {
	// А здесь отправка email
	log.Printf("Отправка письма по событию: %s", string(msg))
	return nil
}
