package ws

import "github.com/go-redis/redis"

type Config struct {
	Redis *redis.Client
}

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	rooms       map[*Room]bool
	redisClient *redis.Client
}

func NewWebsocketHub(c *Config) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		rooms:       make(map[*Room]bool),
		redisClient: c.Redis,
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
