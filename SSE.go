package main

import (
	"encoding/json"
	"fmt"

	"github.com/alexandrevicenzi/go-sse"
)

var SseServer = sse.NewServer(nil)

func PushMessage(message Message) {
	json, err := json.Marshal(message)

	if err != nil {
		fmt.Println(err)
		return
	}
	SseServer.SendMessage("/events/messages", sse.SimpleMessage(string(json)))
}
