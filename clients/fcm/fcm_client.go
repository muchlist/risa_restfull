package fcm

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/muchlist/erru_utils_go/logger"
)

func NewFcmClient() ClientAssumer {
	return &fcmClient{}
}

type fcmClient struct {
}

type ClientAssumer interface {
	SendMessage(payload Payload)
}

// SendMEssage mengirimkan pesan notifikasi ke firebase
func (u *fcmClient) SendMessage(payload Payload) {

	// cek jika penerima 0 atau nil
	if payload.ReceiverTokens == nil || len(payload.ReceiverTokens) == 0 {
		logger.Info("receiver tidak ada")
		return
	}

	// cek jika masing masing penerima tidak string kosong atau nil
	var validTokens []string
	for _, token := range payload.ReceiverTokens {
		if token != "" {
			validTokens = append(validTokens, token)
		}
	}

	// cek jika valid token tidak 0
	if len(validTokens) == 0 {
		logger.Info("receiver tidak ada")
		return
	}

	client, err := FCM.Messaging(context.Background())
	if err != nil {
		logger.Error("gagal mendapatkan client messaging", err)
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Message,
		},
		Tokens: validTokens,
	}

	_, err = client.SendMulticast(context.Background(), message)
	if err != nil {
		logger.Error("gagal mengirim pesan firebase", err)
	}
}
