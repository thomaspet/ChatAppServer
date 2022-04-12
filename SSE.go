package main

import (
	"fmt"

	"github.com/alexandrevicenzi/go-sse"
)

var SseServer = sse.NewServer(nil)

func PushMessage(message Message) {
	// Format message to a "timestamp - author: message" string
	messageString := fmt.Sprintf("%d - %s: %s", message.Timestamp, message.Author, message.Message)
	SseServer.SendMessage("/events/messages", sse.SimpleMessage((messageString)))
}
