package main

import "github.com/ant0ine/go-json-rest/rest"

func Chat(w rest.ResponseWriter, r *rest.Request) {
	var msg Messaging
	err := r.DecodeJsonPayload(&msg)
	if err != nil {
		w.WriteHeader(500)
		w.WriteJson(map[string]string{"error": err.Error()})
		return
	}

	sendTextMessage("", msg.Recipient.ID, "", msg.Message.Text)
}
