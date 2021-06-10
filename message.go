package main

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const BotMessageAction = "bot-message"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  string  `json:"target"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
