package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danyouknowme/smthng/pkg/logger"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingPeriod     = (PongWait * 9) / 10
	MaxMessageSize = 10000
)

var (
	newline  = []byte{'\n'}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	ID    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}

func newClient(conn *websocket.Conn, hub *Hub, id string) *Client {
	return &Client{
		ID:    id,
		conn:  conn,
		hub:   hub,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(MaxMessageSize)

	_ = client.conn.SetReadDeadline(time.Now().Add(PongWait))

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			fmt.Println("send message to client", string(message))
			_ = client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client

	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)

	_ = client.conn.Close()
}

func ServeWs(hub *Hub, ctx *gin.Context) {
	userId := ctx.MustGet(middleware.AuthorizationUserIdKey).(string)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	client := newClient(conn, hub, userId)

	go client.writePump()
	go client.readPump()

	hub.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message domains.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		logger.Errorf("Error on unmarshal JSON message %s", err)
	}

	switch message.Action {
	case "join_channel":
		client.handleJoinChannelMessage(message)
	}
}

func (client *Client) handleJoinRoomMessage(message domains.ReceivedMessage) {
	roomID := message.RoomID

	room := client.hub.findRoomById(roomID)
	if room == nil {
		room = client.hub.createRoom(roomID)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleJoinChannelMessage(message domains.ReceivedMessage) {
	roomID := message.RoomID

	channel, err := client.hub.channelUsecase.GetChannelByID(roomID)
	if err != nil {
		return
	}

	if ok := client.hub.channelUsecase.IsMember(channel.ID, client.ID); !ok {
		return
	}

	client.handleJoinRoomMessage(message)
}
