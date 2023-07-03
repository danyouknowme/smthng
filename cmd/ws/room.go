package ws

import "github.com/go-redis/redis"

type WebsocketMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type Room struct {
	ID         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *WebsocketMessage
	redis      *redis.Client
}
