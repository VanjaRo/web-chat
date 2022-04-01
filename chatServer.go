package main

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
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
	ws.clients[client] = true
}

func (ws *WsServer) unregisterClient(client *Client) {
	delete(ws.clients, client)
}
func (ws *WsServer) broadcastToClients(message []byte) {
	for client := range ws.clients {
		client.send <- message
	}
}
