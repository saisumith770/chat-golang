package RTC

import (
	"net/http"
	"context"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	
	"RTC/sessions"
)

type Server struct {
	RedisState bool
}

var ctx = context.Background()

var Initiation = &Server{true}

func main() {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := redisConn.Ping(ctx).Result()
	if err != nil || pong != "PONG" {
		Initiation.RedisState = false
	}
	router := mux.NewRouter()

	sessions.SessionRouter(
		router.Path("/sessions").Subrouter(),
		redisConn,
		Initiation.RedisState,
	)

	http.ListenAndServe(":8080", router)
}
