package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"github.com/nii236/flux-challenge-go/server/JSONLoader"
)

var darkJedis []JSONLoader.DarkJedi
var worlds []JSONLoader.World

func main() {
	fmt.Println("Starting Flux Challenge server...")
	worlds, darkJedis = JSONLoader.LoadJSON()
	go websocketServer()
	go restServer()
	for {
	}
}

//WebsocketServer starts a websocket server
func restServer() {
	fmt.Println(darkJedis)
	mux := http.NewServeMux()
	mux.HandleFunc("/", restHandler)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func websocketServer() {
	fmt.Println(worlds)
	mux := http.NewServeMux()
	mux.HandleFunc("/", websocketHandler)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":4000")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func restHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received REST request")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received WS request")
	conn, err := upgrader.Upgrade(w, r, nil)
	conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the Flux Challenge!"))
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, p, _ := conn.ReadMessage()
		log.Println(string(p))
	}
}
