package ws

import "github.com/gorilla/websocket"

type Client struct {
	ID    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}
