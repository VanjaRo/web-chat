package main

import (
	"encoding/json"
	"log"

	"github.com/VanjaRo/web-chat/config"
	"github.com/VanjaRo/web-chat/models"
	"github.com/google/uuid"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	users          []models.User
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewWsServer(roomRepository models.RoomRepository, userRepository models.UserRepository) *WsServer {
	var err error
	wsServer := &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}
	wsServer.users, err = userRepository.GetAllUsers()
	if err != nil {
		panic(err)
	}
	return wsServer
}

func (ws *WsServer) Run() {
	go ws.listenPubSubChannel()
	for {
		select {
		case client := <-ws.register:
			ws.registerClient(client)
		case client := <-ws.unregister:
			ws.unregisterClient(client)
		case message := <-ws.broadcast:
			ws.broadcastToClients(message)

		}
	}
}

func (ws *WsServer) registerClient(client *Client) {
	ws.userRepository.AddUser(client)

	ws.publishClientJoined(client)

	ws.listOnlineClients(client)
	ws.clients[client] = true

}

func (ws *WsServer) unregisterClient(client *Client) {
	if _, ok := ws.clients[client]; ok {
		delete(ws.clients, client)

		ws.userRepository.RemoveUser(client)

		ws.publishClientLeft(client)
	}
}
func (ws *WsServer) broadcastToClients(message []byte) {
	for client := range ws.clients {
		client.send <- message
	}
}

func (ws *WsServer) findRoomByName(roomName string) *Room {
	var targetRoom *Room
	for room := range ws.rooms {
		if room.GetName() == roomName {
			targetRoom = room
			break
		}
	}
	if targetRoom == nil {
		targetRoom = ws.runRoomFromRepository(roomName)
	}

	return targetRoom
}

func (ws *WsServer) runRoomFromRepository(roomName string) *Room {
	var room *Room
	dbRoom, err := ws.roomRepository.FindRoomByName(roomName)
	if err == nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.Run()
		ws.rooms[room] = true
	}
	return room
}

func (ws *WsServer) findRoomById(id string) *Room {
	var targetRoom *Room
	for room := range ws.rooms {
		if room.GetId() == id {
			targetRoom = room
			break
		}
	}
	return targetRoom
}

func (ws *WsServer) findClientById(id string) *Client {
	var targetClient *Client
	for client := range ws.clients {
		if client.GetId() == id {
			targetClient = client
			break
		}
	}
	return targetClient
}

func (ws *WsServer) createRoom(roomName string, private bool) *Room {
	room := NewRoom(roomName, private)
	ws.roomRepository.AddRoom(room)

	go room.Run()
	ws.rooms[room] = true
	return room
}

func (ws *WsServer) listOnlineClients(client *Client) {
	for _, user := range ws.users {
		// log.Println("User: ", user.GetName())
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.send <- message.encode()
	}
}

func (ws *WsServer) publishClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}
	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println("Error publishing client joined message: ", err)
	}
}

func (ws *WsServer) publishClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}
	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println("Error publishing client left message: ", err)
	}
}

func (ws *WsServer) listenPubSubChannel() {
	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	ch := pubsub.Channel()
	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Println("Error unmarshalling message: ", err)
			return
		}
		switch message.Action {
		case UserJoinedAction:
			ws.handleUserJoined(message)
		case UserLeftAction:
			ws.handleUserLeft(message)
		case JoinRoomPrivateAction:
			ws.handleUserJoinedPrivate(message)
		}

	}
}

func (ws *WsServer) handleUserJoinedPrivate(message Message) {
	target := ws.findClientById(message.Message)
	if target != nil {
		target.joinRoom(message.Target.GetName(), message.Sender)
	}
}

func (ws *WsServer) handleUserJoined(message Message) {
	ws.users = append(ws.users, message.Sender)
	ws.broadcastToClients(message.encode())
}

func (ws *WsServer) handleUserLeft(message Message) {
	for i, user := range ws.users {
		if user.GetId() == message.Sender.GetId() {
			ws.users[i] = ws.users[len(ws.users)-1]
			ws.users = ws.users[:len(ws.users)-1]
		}
	}
	ws.broadcastToClients(message.encode())
}

func (ws *WsServer) findUserById(ID string) models.User {
	var targetUser models.User
	for _, user := range ws.users {
		if user.GetId() == ID {
			targetUser = user
			break
		}
	}
	return targetUser
}
