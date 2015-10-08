package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("Starting Flux Challenge server...")
	go websocketServer()
	go restServer()
	for {

	}
}

func websocketServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", restHandler)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func restServer() {
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
		messageType, p, err := conn.ReadMessage()
		log.Println(messageType)
		log.Println(string(p))
		log.Println(err)
	}
}
