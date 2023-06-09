package main

import (
	"encoding/json"
	"fmt"
	"log"

	"messaging-service/controllers/controltower"
	"messaging-service/types/entities"
	"messaging-service/types/eventtypes"
	"messaging-service/types/records"
	"messaging-service/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, contentType, Content-Type, Accept, Authorization")
}

func main() {
	r := mux.NewRouter()

	msgController := controltower.New()
	r.HandleFunc("/create-room", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		// todo, extend the 'to' field to be an array
		req := entities.OpenRoomRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}

		// save the room
		roomUUID := uuid.New().String()
		toUUID := utils.ToStr(req.ToUUID)
		fromUUID := utils.ToStr(req.FromUUID)
		room := entities.ChatRoom{
			UUID:         utils.ToStrPtr(roomUUID),
			Participants: []string{toUUID, fromUUID},
		}

		// push this out to the redis server events channel
		openRoomEvent := &entities.OpenRoomEvent{
			FromUUID:  req.FromUUID,
			ToUUID:    req.ToUUID,
			EventType: utils.ToStrPtr(eventtypes.EVENT_OPEN_ROOM.String()),
			Room:      &room,
		}

		newRoom := &records.ChatRoom{
			UUID: roomUUID,
			Participants: []*records.ChatParticipant{
				{
					UUID:     uuid.New().String(),
					RoomUUID: roomUUID,
					UserUUID: fromUUID,
				},
				{
					UUID:     uuid.New().String(),
					RoomUUID: roomUUID,
					UserUUID: toUUID,
				},
			},
		}

		err := msgController.Repo.SaveRoom(newRoom)
		if err != nil {
			panic(err)
		}

		// need to save the room
		msgBytes, err := json.Marshal(openRoomEvent)
		if err != nil {
			panic(err)
		}

		msgController.RedisClient.PublishToRedisChannel(eventtypes.CHANNEL_SERVER_EVENTS, msgBytes)

		w.Write([]byte("created room"))
	}).Methods("POST")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		go msgController.SetupClientConnection(conn)
	})

	r.HandleFunc("/get-rooms-by-user-uuid/{user-uuid}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		vars := mux.Vars(r)
		userUUID, ok := vars["user-uuid"]
		if !ok {
			panic("id is missing in parameters")
		}

		rooms, err := msgController.GetRoomsByUserUUID(userUUID)
		if err != nil {
			panic(err)
		}

		msgController.SubscribeRoomsToServer(rooms, userUUID)

		response := entities.GetRoomsByUserUUIDResponse{
			Rooms: rooms,
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.Write(bytes)
	}).Methods("GET")

	log.Println("Opening server")
	http.ListenAndServe(":9090", r)
}
