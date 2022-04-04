package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/VanjaRo/web-chat/auth"
	"github.com/VanjaRo/web-chat/config"
	"github.com/VanjaRo/web-chat/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name, ID string) *Client {
	client := &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
		Name:     name,
		ID:       uuid.New(),
	}
	client.ID, _ = uuid.Parse(ID)
	return client
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	userCtxVal := r.Context().Value(auth.UserContextKey)
	if userCtxVal == nil {
		log.Printf("Error: User not authenticated")
		return
	}

	user := userCtxVal.(interfaces.User)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}

	client := newClient(conn, wsServer, user.GetName(), user.GetId())

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()
	// concurrent read parameters
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			// if the channel is closed –– close the connection
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// to check if connection is still alive
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		roomId := message.Target.GetId()

		if room := client.wsServer.findRoomById(roomId); room != nil {
			room.broadcast <- &message
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}

}

func (client *Client) joinRoom(roomName string, sender interfaces.User) *Room {
	// sender is the person to start a private chat with
	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		// if sender is not included, it's a public chat
		room = client.wsServer.createRoom(roomName, sender != nil)
	}
	if sender == nil && room.GetPrivate() {
		return nil
	}
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room, sender)
	}
	return room

}

func (client *Client) isInRoom(room *Room) bool {
	_, ok := client.rooms[room]
	return ok
}

func (client *Client) notifyRoomJoined(room *Room, sender interfaces.User) {
	message := Message{
		Action: RoomJoinedAction,
		Sender: sender,
		Target: room,
	}
	client.send <- message.encode()
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message
	client.joinRoom(roomName, nil)
}

func (client *Client) handleJoinRoomPrivateMessage(message Message) {
	target, err := client.wsServer.userRepository.FindUserById(message.Message)
	if err != nil {
		return
	}
	ids := []string{client.GetId(), target.GetId()}
	sort.Strings(ids)
	// create a private room name by concatinating the two client id's
	roomName := strings.Join(ids, "")

	room := client.joinRoom(roomName, target)
	if room != nil {
		client.inviteTargetUser(target, room)
	}

}

func (client *Client) inviteTargetUser(target interfaces.User, room *Room) {
	inviteMessage := Message{
		Action:  JoinRoomPrivateAction,
		Sender:  client,
		Target:  room,
		Message: target.GetId(),
	}
	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Printf("Error publishing message: %v", err)
	}
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	roomId := message.Message
	room := client.wsServer.findRoomById(roomId)
	if room == nil {
		return
	}

	delete(client.rooms, room)
	room.unregister <- client
}

func (client *Client) GetId() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}
