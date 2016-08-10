package main

import (
	"bytes"
	"net/http"
)

func receivedAuthentication(msg string) {}
func receivedMessage(msg Messaging) {
	if msg.Message.Text == "" {
		return
	}

	switch msg.Message.Text {
	// case "image":
	// case "button":
	// case "generic":
	// case "receipt":
	// case "555":
	default:
		sendTextMessage(msg.Sender.ID, msg.Message.Text)
	}
}
func receivedDeliveryConfirmation(msg string) {}
func receivedPostback(msg string)             {}

func sendTextMessage(id, text string) {
	var messageData = `{
		"recipient": {
			"id": "` + id + `"
		},
		"message": {
			"text": "` + text + `"
		}
	}`
	callSendAPI(messageData)
}

func callSendAPI(data string) error {
	c := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	r, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages?access_token="+pageAccessToken, bytes.NewBufferString(data))
	if err != nil {
		return err
	}

	_, err = c.Do(r)
	return err
}
