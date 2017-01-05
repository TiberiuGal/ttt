package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"tibi/lorem/evdispatch/ttt"
)

var moveChan chan ttt.Point
var notify chan struct{}


func Serve(c chan ttt.Point, n chan struct{}) {
	moveChan = c
	notify = n
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/game", getGame).Methods("GET")
	r.HandleFunc("/play", postPlay).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe("localhost:8021", handlers.CORS(headersOk, originsOk, methodsOk)(r))
	if err != nil {
		fmt.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	t.Execute(w, struct{}{})
}

func game(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("game.html")
	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	t.Execute(w, struct{}{})

}

func getGame(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	enc := json.NewEncoder(w)
	err := enc.Encode(game)
	if err != nil {
		fmt.Fprint(w, err)
	}

}

type PlayMsg struct {
	Player int
	Move   Point
}

func postPlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dec := json.NewDecoder(r.Body)
	var p PlayMsg
	dec.Decode(&p)
	fmt.Println("played", p)
	moveChan <- p.Move
	<-notify
	w.Write([]byte("{\"result\":\"ok\"}"))
}
