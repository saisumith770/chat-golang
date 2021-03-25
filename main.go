package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type Server struct {
	RedisState bool
}

var Initiation = &Server{true}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisConn.Ping().Result()
	if err != nil || pong != "PONG" {
		Initiation.RedisState = false
	}

	router := mux.NewRouter()
	hub := NewHub()
	go hub.run()

	router.HandleFunc("/{room}", serveHome)

	router.HandleFunc("/ws/{room}", func(rw http.ResponseWriter, r *http.Request) {
		query := mux.Vars(r)
		ServeWs(hub, rw, r,query["room"])
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
