package notify

type NotificationHandler interface {
	HandleNotification(msg []byte) error
}
