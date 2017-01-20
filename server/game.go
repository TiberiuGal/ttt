package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"tibi/lorem/evdispatch/ttt"
	"github.com/gorilla/mux"
	"github.com/stretchr/objx"
)

type (
	WebGame struct {
		*ttt.Game
		clients map[*client]struct{}
		players map[ttt.Fill]*client
		join    chan *client
		leave   chan *client
		move    chan ttt.Point
		resp 	chan struct{}
	}
	player struct{
		Name interface{}
		Fill ttt.Fill
		Avatar interface{}
	}
	ExportGame struct {
		*ttt.Game
		Players [2]player
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


func initNewGame() int {
	wg := NewGame()
	go wg.run()
	id := len(games)
	games = append(games, wg)
	router.Handle("/wg/"+strconv.Itoa(id), wg)
	return id
}

func newGame(w http.ResponseWriter, r *http.Request) {
	id := initNewGame()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	enc := json.NewEncoder(w)
	res := struct {
		Result string
		Id     int
	}{"ok", id}
	enc.Encode(res)
}

func listGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	res := struct {
		G map[int]struct {
			ttt.Status
			Id int
		}
	}{make(map[int]struct {
		ttt.Status
		Id int
	})}
	for i, g := range games {
		if g.Status != ttt.Finished {
			res.G[i] = struct {
				ttt.Status
				Id int
			}{g.Status, i}
		}
	}
	enc.Encode(res)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	enc := json.NewEncoder(w)
	game := games[id]
	exp := ExportGame{Game:game.Game}
	if p, ok := game.players[ttt.FillX]; ok {
		exp.Players[0] = player{p.userData["name"], p.fill, p.userData["avatar"]}
	}
	if p, ok := game.players[ttt.FillO]; ok {
		exp.Players[1] = player{p.userData["name"], p.fill, p.userData["avatar"]}
	}

	err := enc.Encode(exp)
	if err != nil {
		fmt.Fprint(w, err)
	}

}


func NewGame() *WebGame {
	wg := &WebGame{
		Game:    ttt.NewGame(),
		clients: make(map[*client]struct{}),
		players: make(map[ttt.Fill]*client),
		join: make(chan *client),
		leave: make(chan *client),
		move: make(chan ttt.Point),
		resp: make(chan struct{}),
	}
	return wg
}

func (g *WebGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving http reuqest on webgame")

	socket, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("socket upgraded")
	if err != nil {
		log.Fatal("ServeHTTP webgame cannot upgrade socket", err)
		return
	}

	authCookie, err := r.Cookie("auth")
	if err != nil {
		log.Fatal("failed to get auth cokie ", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan string, messageBufferSize),
		room:     g,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	n := len(g.players)
	rejoined := false
	for f, c:= range g.players {
		if c.userData["email"] == client.userData["email"] {
			fmt.Println("rejoining client", client.userData["email"])
			client.fill = f
			g.players[f] = client
			rejoined = true
			break

		}
	}

	if ! rejoined {
		log.Println("new client is joining")
		if n == 0 {
			client.fill = ttt.FillX
			g.players[ttt.FillX] = client
		} else if n == 1 {
			client.fill = ttt.FillO
			g.players[ttt.FillO] = client
		}
	}

	g.join <- client
	log.Println("new client joined the room")

	defer func() { g.leave <- client }()
	go client.write()
	j, _ := json.Marshal(struct{Result string; Action string; Fill ttt.Fill}{"ok", "join", client.fill})
	client.send <- string(j)
	fmt.Println("we have", len(g.players), "joined, status is", g.Status )
	if len(g.players) == 2  && g.Status == ttt.Waiting {
		g.Status = ttt.Ready
	}
	if g.Status == ttt.Ready  {

		for c := range g.clients {
			j, _ := json.Marshal(struct{Result string; Action string; Fill ttt.Fill}{"ok", "start", g.Next})
			c.send <- string(j)
		}
	}
	client.read()
}

func (g *WebGame) run() {
	log.Print("start running")
	go g.GameLoop(g.move, g.resp)
	for {
		select {
		case client := <-g.join:
			log.Print("received join request")
			g.clients[client] = struct{}{}

		case client := <-g.leave:
			log.Print("received leave")
			delete(g.clients, client)
			close(client.send)


		case <-g.resp:
			log.Print("received move done ")
			j, _ := json.Marshal(struct{Result string; Action string; Fill ttt.Fill}{"ok", "move", g.Next})
			for client := range g.clients {

				select {
				case client.send <- string(j):
				default:
					delete(g.clients, client)
					close(client.send)

				}
			}

		}
	}
	log.Print("end running")
}
