package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	newline = []byte("hello")
	// space   = []byte{' '}
	timeStamp = time.Now().Format("3:04 PM")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	Name     string `json:"name"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
		Name:     name,
	}

}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		var message Message

		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			fmt.Println("error unmarshal")
		}

		message.Time = timeStamp

		fmt.Println("pesan: ", message)

		a := []rune(message.Message)

		isCode := string(a[0:1])

		if isCode == "/" {
			code := string(a[1:])
			code = getMessageCode(code)
			if code == "" {
				code = "Perintah yang anda kirim tidak ada"
			}
			message.Message = code

			msg, _ := json.Marshal(message)
			client.handleNewMessage(msg)
			if message.Action == SendMessageAction || message.Action == BotMessageAction {
				insertToDb(message)
			}
		} else {
			client.handleNewMessage(jsonMessage)
			if message.Action == SendMessageAction {
				insertToDb(message)
			}
		}

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
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err != nil {
				log.Fatal(err)
			}

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
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

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message
	room := client.wsServer.findRoomByName(roomName)

	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	count := countMsg()
	if count != 0 {
		a := getPreviousMsg(roomName)
		for _, msg := range a {
			client.conn.WriteJSON(msg)
		}
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByName(message.Message)
	delete(client.rooms, room)

	room.unregister <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = client
	message.Time = timeStamp

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomName := message.Target
		// Use the ChatServer method to find the room, and if found, broadcast!
		room := client.wsServer.findRoomByName(roomName)

		if room != nil {
			room.broadcastToClientsInRoom(message.encode())
		}

	// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	case BotMessageAction:
		roomName := message.Target
		room := client.wsServer.findRoomByName(roomName)

		messageBot := &Message{
			Action:  SendMessageAction,
			Target:  room.Name,
			Message: message.Message,
			Sender: &Client{
				conn:     client.conn,
				wsServer: client.wsServer,
				send:     client.send,
				rooms:    client.rooms,
				Name:     "Topin",
			},
		}

		room.broadcastToClientsInRoom(messageBot.encode())
	}
}
