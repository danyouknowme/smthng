package ws

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type Room struct {
	ID          string
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *domains.WebsocketMessage
	redisClient *redis.Client
}

var ctx = context.Background()

func NewRoom(id string, rds *redis.Client) *Room {
	return &Room{
		ID:          id,
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *domains.WebsocketMessage),
		redisClient: rds,
	}
}

func (room *Room) RunRoom() {
	go room.subscribeToRoomMessages()

	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.publishRoomMessage(message.Encode())
		}
	}
}

func (room *Room) GetId() string {
	return room.ID
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	delete(room.clients, client)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) publishRoomMessage(message []byte) {
	err := room.redisClient.Publish(ctx, room.GetId(), message).Err()

	if err != nil {
		logger.Error(err)
	}
}

func (room *Room) subscribeToRoomMessages() {
	pubsub := room.redisClient.Subscribe(ctx, room.GetId())

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}
