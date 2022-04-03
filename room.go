package main

import "github.com/google/uuid"

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
	for {
		select {
		case client := <-room.register:
			room.registerClient(client)
		case client := <-room.unregister:
			room.unregisterClient(client)
		case message := <-room.broadcast:
			room.broadcastToClients(message.encode())

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
	room.broadcastToClients(message.encode())
}

func (room *Room) isPrivate() bool {
	return room.Private
}
