package main

import (
	"context"
	"encoding/json"
	"log"

	"messaging-service/controllers/controltower"
	"messaging-service/handlers"
	"messaging-service/repo"
	"net/http"

	redisClient "messaging-service/redis"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, contentType, Content-Type, Accept, Authorization")
}

type rootHandler func(http.ResponseWriter, *http.Request) (interface{}, error)

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	res, err := fn(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// https://markphelps.me/posts/handling-errors-in-your-http-handlers/
func main() {
	fx.New(
		fx.Provide(NewMuxRouter),
		fx.Provide(handlers.New),
		fx.Provide(redisClient.New),
		fx.Provide(controltower.New),
		fx.Provide(repo.New),
		fx.Invoke(SetupRoutes),
	).Run()
}

func SetupRoutes(h *handlers.Handler, r *mux.Router) {

	// websocket
	r.HandleFunc("/ws", h.SetupWebsocketConnection)

	// API
	r.Handle("/delete-room", rootHandler(h.DeleteRoom)).Methods("POST")
	r.Handle("/leave-room", rootHandler(h.LeaveRoom)).Methods("POST")
	r.Handle("/create-room", rootHandler(h.CreateRoom)).Methods("POST")
	r.Handle("/get-messages-by-room-uuid", rootHandler(h.GetMessagesByRoomUUID)).Methods("GET")
	r.Handle("/get-rooms-by-user-uuid", rootHandler(h.GetRoomsByUserUUID)).Methods("GET")

	h.SetupChannels()
}

func NewMuxRouter(lc fx.Lifecycle) *mux.Router {
	r := mux.NewRouter()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Opening server on port 9090")
				err := http.ListenAndServe(":9090", r)
				if err != nil {
					log.Println(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return r
}
