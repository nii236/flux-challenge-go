package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"github.com/nii236/flux-challenge-go/server/JSONLoader"
)

var darkJedis []JSONLoader.DarkJedi
var worlds []JSONLoader.World
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}
var baseURL = "localhost"
var restPort = ":3000"
var wsPort = ":4000"

//AugmentedJedi will return Dark Jedi struct with Base URL
type AugmentedJedi func(jedi JSONLoader.DarkJedi) JSONLoader.AugmentedDarkJedi

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("Starting Flux Challenge server...")
	worlds, darkJedis = JSONLoader.LoadJSON()
	go websocketServer()
	go restServer()
	for {
	}
}

func restServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dark-jedis", restHandler)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(baseURL + restPort)

}

func restHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received REST request")
	query := r.URL.Query()
	darkJediID, _ := strconv.Atoi(query.Get("id"))
	log.Println(darkJediID)
	if darkJediID != 0 {
		sendDarkjedi(w, r)
	} else {
		sendRandomDarkJedi(w, r)
	}

}

func sendDarkjedi(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	darkJediID, _ := strconv.Atoi(query.Get("id"))
	for _, sith := range darkJedis {
		if darkJediID == sith.ID {
			// time.Sleep(rand.Intn(1000) * time.Millisecond)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			augmentedDJ := augmentJediWithNeighbourURLs()(sith)
			json.NewEncoder(w).Encode(augmentedDJ)
			return
		}
	}
	fmt.Fprint(w, "Dark Jedi not found")
}

func sendRandomDarkJedi(w http.ResponseWriter, r *http.Request) {

	sith := darkJedis[rand.Intn(len(darkJedis)-1)]
	augmentedSith := augmentJediWithNeighbourURLs()(sith)
	json.NewEncoder(w).Encode(augmentedSith)
}

func websocketServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", websocketHandler)
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(baseURL + wsPort)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received WS request")
	conn, err := upgrader.Upgrade(w, r, nil)
	worldChan := make(chan JSONLoader.World)
	go sendRandomWorld(worldChan)
	if err != nil {
		log.Fatal(err)
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the Flux Challenge!"))
	if err != nil {
		log.Println(err)
		return
	}
	for {
		select {
		case w := <-worldChan:
			j, err := json.Marshal(w)
			if err != nil {
				log.Fatal(err)
			}
			conn.WriteMessage(websocket.TextMessage, []byte(j))
		}
	}
}

func sendRandomWorld(worldChan chan JSONLoader.World) {
	fmt.Println("Begin random world sender")
	for {
		time.Sleep(5 * time.Second)
		w := worlds[rand.Intn(len(worlds)-1)]
		worldChan <- w
	}
}

func augmentJediWithNeighbourURLs() AugmentedJedi {
	return func(jedi JSONLoader.DarkJedi) JSONLoader.AugmentedDarkJedi {
		dj := JSONLoader.AugmentedDarkJedi{
			ID:        jedi.ID,
			Name:      jedi.Name,
			Homeworld: jedi.Homeworld,
		}
		if &dj.Master != nil {
			dj.Master.URL = fmt.Sprint(baseURL, restPort, "/dark-jedis?id=", jedi.Master)
			dj.Master.ID = jedi.Master
			dj.Apprentice.URL = fmt.Sprint(baseURL, restPort, "/dark-jedis?id=", jedi.Apprentice)
			dj.Apprentice.ID = jedi.Apprentice
		}
		log.Println("Augmented Jedi struct:", dj)
		return dj
	}

}
