package sessions

import (
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func SessionRouter(
	router *mux.Router,
	client *redis.Client,
	RedisState bool,
) {
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				rw.WriteHeader(http.StatusExpectationFailed)
			}
		}()
		if RedisState {
			body, BodyErr := ioutil.ReadAll(r.Body)
			if BodyErr != nil {
				rw.WriteHeader(http.StatusBadRequest)
			} else {
				sessionID := Create_new_session(client, body)
				_, err := rw.Write([]byte(sessionID))
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}).Methods("POST")

	router.HandleFunc("/{session_id}", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				rw.WriteHeader(http.StatusConflict)
			}
		}()
		queries := mux.Vars(r)
		if RedisState {
			Erase_session(client, queries["session_id"])
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}).Methods("DELETE")

	router.HandleFunc("/{session_id}", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				rw.WriteHeader(http.StatusConflict)
			}
		}()
		queries := mux.Vars(r)
		if RedisState {
			rw.Write(Get_session(client, queries["session_id"]))
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}).Methods("GET")
}
