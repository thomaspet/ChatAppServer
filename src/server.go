package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocraft/web"
	"github.com/google/uuid"
)

type HttpContext struct {
	UserName string
	Guid     string
}
type Message struct {
	Message    string
	Timestamp  string
	Author     string
	AuthorGuid string
}

type MessagesResponse struct {
	Messages []Message
	Token    string
}

var messages []Message = make([]Message, 100)
var shortLivedTokens map[string]int64 = make(map[string]int64)

func (ctx *HttpContext) SendMessage(rw web.ResponseWriter, req *web.Request) {
	bytes, _ := ioutil.ReadAll(req.Body)

	var message Message = Message{
		Message:    string(bytes),
		Timestamp:  time.Now().Format(time.RFC3339),
		Author:     ctx.UserName,
		AuthorGuid: ctx.Guid,
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
		if message.Timestamp != "" {
			createdMessages = append(createdMessages, message)
		}
	}

	token := uuid.New().String()
	var messagesResponse MessagesResponse = MessagesResponse{
		Messages: createdMessages,
		Token:    token,
	}

	shortLivedTokens[token] = time.Now().Unix() + 60 // 1 minute TTL

	jsonbody, err := json.MarshalIndent(messagesResponse, "", "  ")
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(jsonbody))

}
