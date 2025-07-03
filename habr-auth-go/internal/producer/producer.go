package producer

type EventProducer interface {
	SendRegisterEvent(email string, code int) error
}
