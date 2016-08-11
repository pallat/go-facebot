package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

var c = http.DefaultClient

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

func receivedPostback(msg Messaging) {
	sendTextMessage("", msg.Sender.ID, "", "Postback called")
}

func sendTextMessage(pageID, id, mid, text string) error {
	var messageData = `{
  "recipient":{
    "id":"` + id + `"
  },
  "message":{
    "text":"` + text + `"
  }`

	// "sender":{
	//   "id":"` + pageID + `"
	// },
	//             "mid":"` + mid + `",

	return callSendAPI(messageData)
}

func callSendAPI(data string) error {
	// c := &http.Client{
	// 	Transport: &http.Transport{
	// 	// DisableKeepAlives: true,
	// 	},
	// }

	r, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages?access_token="+pageAccessToken, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	br, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(br))

	resp, err := c.Do(r)

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))

	return err
}
