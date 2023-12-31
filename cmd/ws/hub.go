package ws

import (
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	ChannelUsecase usecases.ChannelUsecase
	Redis          *redis.Client
}

type Hub struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	channelUsecase usecases.ChannelUsecase
	redisClient    *redis.Client
}

func NewWebsocketHub(c *Config) *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		channelUsecase: c.ChannelUsecase,
		redisClient:    c.Redis,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.registerClient(client)

		case client := <-hub.unregister:
			hub.unregisterClient(client)

		case message := <-hub.broadcast:
			hub.broadcastToClients(message)
		}
	}
}

func (hub *Hub) registerClient(client *Client) {
	hub.clients[client] = true
}

func (hub *Hub) unregisterClient(client *Client) {
	delete(hub.clients, client)
}

func (hub *Hub) broadcastToClients(message []byte) {
	for client := range hub.clients {
		client.send <- message
	}
}

func (hub *Hub) BroadcastToRoom(message []byte, roomID string) {
	if room := hub.findRoomByID(roomID); room != nil {
		room.publishRoomMessage(message)
	}
}

func (hub *Hub) findRoomByID(id string) *Room {
	var foundRoom *Room
	for room := range hub.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (hub *Hub) createRoom(id string) *Room {
	room := NewRoom(id, hub.redisClient)
	go room.RunRoom()
	hub.rooms[room] = true

	return room
}
