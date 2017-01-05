package server

import (
	"log"

	"github.com/gorilla/websocket"
	"tibi/lorem/evdispatch/ttt"
)

type (

	client struct {

		// socket is the web socket for this client.
		socket *websocket.Conn

		// send is a channel on which messages are sent.
		send chan ttt.Point

		room *WebGame
	}
)



func (c *client) read() {
	for {
		var msg ttt.Point
		if err := c.socket.ReadJSON(&msg); err == nil {
			c.room.move <- msg
		} else {
			log.Println("invalid json read from client", err)
			break

		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}