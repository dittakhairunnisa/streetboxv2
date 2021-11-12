package fcm

import "firebase.google.com/go/v4/messaging"

// ServiceInterface ..
type ServiceInterface interface {
	Unsubscribe(topic string, token []string) error
	SendNotification(
		topic string,
		notif *messaging.Notification) string
	SendNotificationWithData(
		topic string,
		notif *messaging.Notification,
		data *map[string]string,
		channelID string) string
	SendNotificationWithDataToken(
		token string,
		notif *messaging.Notification,
		data *map[string]string,
		channelID string) string
	SendNotificationToOne(
		notif *messaging.Notification,
		token string) string
}
