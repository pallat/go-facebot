package main

type Common struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID         string      `json:"id"`
	Time       int64       `json:"time"`
	Messagings []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    IDs      `json:"sender"`
	Recipient IDs      `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   *Message `json:"message,omitempty"`
}

type IDs struct {
	ID string `json:"id"`
}

type Message struct {
	MID        string     `json:"mid"`
	SEQ        int        `json:"seq"`
	Text       string     `json:"text"`
	QuickReply QuickReply `json:"quick_reply"`
}

type QuickReply struct {
	Payload string `json:"payload"`
}
