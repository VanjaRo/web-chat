package main

import (
	"encoding/json"
	"log"
)

// types of actions for a message
const (
	SendMessageAction     = "send-message"
	JoinRoomAction        = "join-room"
	JoinRoomPrivateAction = "join-room-private"
	LeaveRoomAction       = "leave-room"
	UserJoinedAction      = "user-join"
	UserLeftAction        = "user-left"
	RoomJoinedAction      = "room-joined"
)

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  *Room   `json:"target"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}
