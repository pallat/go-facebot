package main

import (
	"bytes"
	"net/http"
)

func receivedAuthentication(msg Messaging) error {
	return sendTextMessage("", msg.Sender.ID, "", "Authentication successful")
}
func receivedMessage(pageID string, msg Messaging) error {
	if msg.Message.Text == "" {
		return nil
	}

	switch msg.Message.Text {
	// case "image":
	// case "button":
	// case "generic":
	// case "receipt":
	// case "555":
	default:
		return sendTextMessage(pageID, msg.Sender.ID, msg.Message.MID, msg.Message.Text)
	}

	return nil
}
func receivedDeliveryConfirmation(msg string) {}
func receivedPostback(msg string)             {}

func sendTextMessage(pageID, id, mid, text string) error {
	var messageData = `{
		"recipient": {
			"id": "` + id + `"
		},
		"message": {
			"text": "` + text + `"
		}
	}`

	// "sender":{
	//   "id":"` + pageID + `"
	// },
	//             "mid":"` + mid + `",

	return callSendAPI(messageData)
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
