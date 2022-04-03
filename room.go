package main

import (
	"context"
	"log"

	"github.com/VanjaRo/web-chat/config"
	"github.com/google/uuid"
)

var ctx = context.Background()

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Private    bool      `json:"private"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(name string, private bool) *Room {
	return &Room{
		Name:       name,
		Private:    private,
		ID:         uuid.New(),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) Run() {
	go room.subscribeRoomMessage()
	for {
		select {
		case client := <-room.register:
			room.registerClient(client)
		case client := <-room.unregister:
			room.unregisterClient(client)
		case message := <-room.broadcast:
			room.publishRoomMessage(message.encode())

		}
	}
}

func (room *Room) registerClient(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClient(client *Client) {
	delete(room.clients, client)
}

func (room *Room) broadcastToClients(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: client.GetName() + " joined the room",
	}
	room.publishRoomMessage(message.encode())
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

func (room *Room) publishRoomMessage(message []byte) {
	if err := config.Redis.Publish(ctx, room.GetName(), message).Err(); err != nil {
		log.Println("Error publishing message to room:", err)
	}
}

func (room *Room) subscribeRoomMessage() {
	pubsub := config.Redis.Subscribe(ctx, room.GetName())
	ch := pubsub.Channel()
	for msg := range ch {
		room.broadcastToClients([]byte(msg.Payload))
	}
}
