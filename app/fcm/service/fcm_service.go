package service

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/errorutils"
	"firebase.google.com/go/v4/messaging"
	"streetbox.id/app/fcm"
)

// FCMService fcm send push notif etc
type FCMService struct {
	App *firebase.App
	Ctx *context.Context
}

// New ..
func New(ctx *context.Context, app *firebase.App) fcm.ServiceInterface {
	return &FCMService{app, ctx}
}

// SendMessage ...
func (s *FCMService) send(message *messaging.Message, topic string) string {
	client, err := s.App.Messaging(context.Background())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	resp, err := client.Send(context.Background(), message)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	return resp
}

// SendNotification ..
func (s *FCMService) SendNotification(
	topic string,
	notif *messaging.Notification) string {
	androidConfig := &messaging.AndroidConfig{
		Priority: "high",
	}
	msg := &messaging.Message{Notification: notif, Android: androidConfig, Topic: topic}
	return s.send(msg, topic)
}

// Unsubscribe after client receive fcm
func (s *FCMService) Unsubscribe(topic string, token []string) error {
	if _, err := s.isSubscribe(false, token, topic); err != nil {
		return err
	}
	return nil
}

// SendNotificationWithData ..
func (s *FCMService) SendNotificationWithData(
	topic string,
	notif *messaging.Notification,
	data *map[string]string,
	channelID string) string {
	androidConfig := &messaging.AndroidConfig{
		Priority: "high",
		Notification: &messaging.AndroidNotification{
			Title: notif.Title,
			Body: notif.Body,
			// ImageURL: notif.ImageURL,
			ChannelID: channelID,
		},
	}
	// msg := &messaging.Message{Notification: notif, Android: androidConfig, Data: *data, Topic: topic}
	msg := &messaging.Message{Android: androidConfig, Data: *data, Topic: topic}
	return s.send(msg, topic)
}

// SendNotificationWithDataToken ..
func (s *FCMService) SendNotificationWithDataToken(
	token string,
	notif *messaging.Notification,
	data *map[string]string,
	channelID string) string {
	client, err := s.App.Messaging(*s.Ctx)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	androidConfig := &messaging.AndroidConfig{
		Priority: "high",
		Notification: &messaging.AndroidNotification{
			Title: notif.Title,
			Body: notif.Body,
			// ImageURL: notif.ImageURL,
			ChannelID: channelID,
		},
	}
	// msg := &messaging.Message{Notification: notif, Android: androidConfig, Data: *data, Token: token}
	msg := &messaging.Message{Android: androidConfig, Data: *data, Token: token}
	resp, err := client.Send(*s.Ctx, msg)
	if messaging.IsUnregistered(err) {
		log.Printf("Registration Token %s invalid", token)
	}
	if messaging.IsSenderIDMismatch(err) {
		log.Print("invalid credential or permission error")
	}
	httpResp := errorutils.HTTPResponse(err)
	log.Printf("ERROR: http response -> %+v", httpResp)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	return resp
}

func (s *FCMService) isSubscribe(isSubs bool, token []string, topic string) (*messaging.Client, error) {
	client, err := s.App.Messaging(context.Background())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	resp := new(messaging.TopicManagementResponse)
	if isSubs {
		resp, err = client.SubscribeToTopic(context.Background(), token, topic)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return nil, err
		}
		log.Printf("INFO: %d tokens were subscribed successfully : %s",
			resp.SuccessCount, token)
		if len(resp.Errors) > 0 {
			for _, v := range resp.Errors {
				log.Printf("INFO: Error Info: %+v", v)
			}
		}

		return client, nil
	}
	resp, err = client.UnsubscribeFromTopic(context.Background(), token, topic)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	log.Printf("INFO: %d tokens were unsubscribed successfully : %s",
		resp.SuccessCount, token)
	return client, nil
}

// SendNotificationToOne ..
func (s *FCMService) SendNotificationToOne(
	notif *messaging.Notification,
	token string) string {
	client, err := s.App.Messaging(*s.Ctx)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	androidConfig := &messaging.AndroidConfig{
		Priority: "high",
	}
	msg := &messaging.Message{Notification: notif, Android: androidConfig, Token: token}
	resp, err := client.Send(*s.Ctx, msg)
	if messaging.IsUnregistered(err) {
		log.Printf("Registration Token %s invalid", token)
	}
	if messaging.IsSenderIDMismatch(err) {
		log.Print("invalid credential or permission error")
	}
	httpResp := errorutils.HTTPResponse(err)
	log.Printf("ERROR: http response -> %+v", httpResp)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err.Error()
	}
	return resp
}
