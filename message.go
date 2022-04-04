package main

import (
	"encoding/json"
	"log"

	"github.com/VanjaRo/web-chat/interfaces"
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
	Action  string          `json:"action"`
	Message string          `json:"message"`
	Target  *Room           `json:"target"`
	Sender  interfaces.User `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
	}
	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
