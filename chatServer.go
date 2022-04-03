package main

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

func (ws *WsServer) Run() {
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
	ws.notifyClientJoined(client)
	ws.listOnlineClients(client)
	ws.clients[client] = true
}

func (ws *WsServer) unregisterClient(client *Client) {
	if _, ok := ws.clients[client]; ok {
		delete(ws.clients, client)
		ws.notifyClientLeft(client)
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
	return targetRoom
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
	go room.Run()
	ws.rooms[room] = true
	return room
}

func (ws *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}
	ws.broadcastToClients(message.encode())
}

func (ws *WsServer) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}
	ws.broadcastToClients(message.encode())
}

func (ws *WsServer) listOnlineClients(client *Client) {
	for onlineClient := range ws.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: onlineClient,
		}
		client.send <- message.encode()
	}
}
