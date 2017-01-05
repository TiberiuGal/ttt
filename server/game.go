package server

import (
	"tibi/lorem/evdispatch/ttt"
	"github.com/gorilla/websocket"

	"log"

	"net/http"
)

type (
	WebGame struct{
		*ttt.Game
		clients map[*client]struct{}
		join chan *client
		leave chan *client
		move chan ttt.Point
	}

)
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var (
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)


func NewGame() *WebGame {
	return &WebGame{
		Game: ttt.NewGame(),
		clients: make(map[*client]struct{}, 0),
	}
}


func (g *WebGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP webgame cannot upgrade socket", err)
		return
	}


	client := &client{
		socket:   socket,
		send:     make(chan ttt.Point, messageBufferSize),
		room : g,

	}

	g.join <- client
	defer func() { g.leave <- client }()
	go client.write()
	client.read()
}

func (g *WebGame) run() {
	for {
		select {
		case client := <-g.join:
			g.clients[client] = struct{}{}

		case client := <-g.leave:
			delete(g.clients, client)
			close(client.send)

		case msg := <-g.move:


			for client := range g.clients {
				select {
				case client.send <- msg:
					//do nothing?

				default:
					delete(g.clients, client)
					close(client.send)

				}
			}
		}
	}
}