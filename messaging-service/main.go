package main

import (
	"encoding/json"
	"log"

	"messaging-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	h := handlers.New()

	// websocket
	r.HandleFunc("/ws", h.SetupWebsocketConnection)

	// API
	r.Handle("/delete-room", rootHandler(h.DeleteRoom)).Methods("POST")
	r.Handle("/create-room", rootHandler(h.CreateRoom)).Methods("POST")
	r.Handle("/get-messages-by-room-uuid", rootHandler(h.GetMessagesByRoomUUID)).Methods("GET")
	r.Handle("/get-rooms-by-user-uuid", rootHandler(h.GetRoomsByUserUUID)).Methods("GET")

	log.Println("Opening server...")
	http.ListenAndServe(":9090", r)
}
