package server

import (
	"fmt"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	games  []*WebGame
	router *mux.Router
)


func Serve(addr string) {
	r := mux.NewRouter()
	router = r
	initNewGame()
	initAuth(addr)
	r.Handle("/", MustAuth(&templateHandler{filename:"index.html"}))
	r.HandleFunc("/auth/{action}", loginHandler)

	r.HandleFunc("/game/{id}", getGame).Methods("GET")
	//r.HandleFunc("/play", postPlay).Methods("POST")
	r.HandleFunc("/newGame", newGame).Methods("POST")
	r.HandleFunc("/listGames", listGames).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir( path.Join("..","static")))))
	err := http.ListenAndServe(addr, handlers.CORS(headersOk, originsOk, methodsOk)(r))
	if err != nil {
		fmt.Println(err)
	}
}






