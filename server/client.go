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
		send chan string
		fill ttt.Fill
		room *WebGame
		userData map[string]interface{}
	}
)



func (c *client) read() {
	log.Println("start reading")
	for {
		var msg ttt.Point
		m := struct{Message ttt.Point
			Txt string}{}

		if err := c.socket.ReadJSON(&m); err != nil {
			log.Println("invalid json read from client", err)
			break
		}
		msg = m.Message
		log.Println("new move",  m , msg)
		c.room.move <- msg

	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(struct{Txt string}{msg}); err != nil {
			break
		}
	}
	c.socket.Close()
}