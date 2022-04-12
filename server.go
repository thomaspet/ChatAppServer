package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

type HttpContext struct{}
type Message struct {
	Message    string
	Timestamp  int64
	Author     string
	AuthorGuid string
}

var messages []Message = make([]Message, 100)

func (ctx *HttpContext) SendMessage(rw web.ResponseWriter, req *web.Request) {
	bytes, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(bytes))

	var message Message = Message{
		Message:    string(bytes),
		Timestamp:  time.Now().Unix(),
		Author:     "Anonymous",
		AuthorGuid: "",
	}

	// Add to messages, if it's full remove the oldest message
	if len(messages) == cap(messages) {
		messages = append(messages[1:], message)
	} else {
		messages = append(messages, message)
	}

	PushMessage(message)

	rw.WriteHeader(http.StatusOK)
}

func (ctx *HttpContext) GetMessages(rw web.ResponseWriter, req *web.Request) {
	// Return messages as json
	createdMessages := []Message{}
	for _, message := range messages {
		if message.Timestamp != 0 {
			createdMessages = append(createdMessages, message)
		}
	}

	jsonbody, err := json.MarshalIndent(createdMessages, "", "  ")
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(jsonbody))

}
