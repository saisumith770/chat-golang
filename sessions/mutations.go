package sessions

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/rogpeppe/fastuuid"
)

type userSession struct {
	user_id  string `json:"user_id"`
	username string `json:"username"`
	photo    string `json:"photo"`
	domain   string `json:"domain"`
}

var ctx = context.Background()

func Create_new_session(
	Conn *redis.Client, 
	payload []byte,
) string {
	var session = fastuuid.Hex128((fastuuid.MustNewGenerator().Next()))
	err := Conn.Set(ctx,session, payload, 0).Err()
	if err != nil {
		panic(err)
	}
	return session
}

func Get_session(
	Conn *redis.Client, 
	session_id string,
) []byte {
	val, err := Conn.Get(ctx,session_id).Result()
	if err != nil {
		panic(err)
	}
	return []byte(val)
}

func Erase_session(
	Conn *redis.Client, 
	session_id string,
) {
	err := Conn.Del(ctx,session_id).Err()
	if err != nil {
		panic(err)
	}
}
