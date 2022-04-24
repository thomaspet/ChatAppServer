package main

import (
	"encoding/json"
	"log"

	"github.com/alexandrevicenzi/go-sse"
)

var SseServer = sse.NewServer(nil)

func PushMessage(message Message) {
	json, err := json.Marshal(message)

	if err != nil {
		log.Panicln(err)
		return
	}
	SseServer.SendMessage("/events/messages", sse.SimpleMessage(string(json)))
}
