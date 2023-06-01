package main

import (
	"encoding/json"
	"errors"
	"fmt"

	redisClient "chatapi/redis"
	"chatapi/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type MessageType int64

const (
	EVENT_CHAT_TEXT MessageType = iota
	EVENT_CHAT_TEXT_METADATA
	EVENT_OPEN_ROOM         // open a chat room request
	EVENT_SET_CLIENT_SOCKET // set the client socket
)

const (

	// server channel for server side events
	CHANNEL_SERVER_EVENTS = "CHANNEL_SERVER_EVENTS"
)

func (m MessageType) String() string {
	switch m {
	case EVENT_CHAT_TEXT:
		return "EVENT_CHAT_TEXT"
	case EVENT_CHAT_TEXT_METADATA:
		return "EVENT_CHAT_TEXT_METADATA"
	case EVENT_OPEN_ROOM:
		return "EVENT_OPEN_ROOM"
	case EVENT_SET_CLIENT_SOCKET:
		return "EVENT_SET_CLIENT_SOCKET"
	}
	return "UNKNOWN"
}

type MessageController struct {
	RedisClient *redisClient.RedisClient

	Connections map[string]*types.Connection

	// map the user uuid to a list of the user's connections (different devices)
	UserConnections         map[string]map[string]bool
	OutboundMessagesChannel <-chan *redis.Message

	// track active rooms/channels on this server
	ActiveChannels map[string]*types.Channel
}

func (c *MessageController) handleIncomingTextMessageFromRedis(msg string) error {
	chatMessage := types.ChatMessageEvent{}
	err := json.Unmarshal([]byte(msg), &chatMessage)
	if err != nil {
		panic(err)
	}

	room, ok := c.ActiveChannels[chatMessage.RoomUUID]
	if !ok {
		return nil
	}

	// get all the outbound connections we need to send the message
	outboundConnections := []*types.Connection{}
	for participantUUID, _ := range room.ParticipantsOnServer {

		connections := c.UserConnections[participantUUID]
		for connUUID := range connections {
			if connUUID != chatMessage.FromConnectionUUID {
				connection, ok := c.Connections[connUUID]
				if !ok {
					continue
				}
				outboundConnections = append(outboundConnections, connection)
			}
		}
	}

	for _, outboundConn := range outboundConnections {
		outboundConn.Conn.WriteJSON(chatMessage)
	}

	return nil
}

func getEventType(event string) (string, error) {
	e := map[string]interface{}{}
	err := json.Unmarshal([]byte(event), &e)
	if err != nil {
		return "", err
	}

	eType, ok := e["eventType"]
	if !ok {
		return "", errors.New("no event type present")
	}
	val, ok := eType.(string)
	if !ok {
		return "", errors.New("could not cast to event type")
	}
	return val, nil
}

func (c *MessageController) handleIncomingServerEventFromRedis(event string) error {
	eventType, err := getEventType(event)
	if err != nil {
		panic(err)
	}

	if eventType == EVENT_OPEN_ROOM.String() {
		openRoomEvent := types.OpenRoomEvent{}
		err = json.Unmarshal([]byte(event), &openRoomEvent)
		if err != nil {
			panic(err)
		}

		var listOfFromConnections, listOfToConnections map[string]bool
		_listOfFromConnections, ok := c.UserConnections[openRoomEvent.FromUUID]
		if ok {
			listOfFromConnections = _listOfFromConnections
		}
		_listOfToConnections, ok := c.UserConnections[openRoomEvent.FromUUID]
		if ok {
			listOfToConnections = _listOfToConnections
		}

		for connUUID := range listOfFromConnections {
			channel, ok := c.ActiveChannels[openRoomEvent.Room.UUID]
			if !ok {
				roomSubscriber := c.RedisClient.SetupChannel(openRoomEvent.Room.UUID)
				go c.subscribeToRedisChannel(roomSubscriber, c.handleIncomingTextMessageFromRedis)

				channel = &types.Channel{
					Subscriber: roomSubscriber,
					UUID:       openRoomEvent.Room.UUID,
				}
				c.ActiveChannels[openRoomEvent.Room.UUID] = channel
			}
			channel.ParticipantsOnServer[openRoomEvent.FromUUID] = true
			c.Connections[connUUID].Conn.WriteJSON(openRoomEvent)
		}

		for connUUID := range listOfToConnections {
			// TODO – use redis client to check if channel is already subscribed
			channel, ok := c.ActiveChannels[openRoomEvent.Room.UUID]
			if !ok {
				roomSubscriber := c.RedisClient.SetupChannel(openRoomEvent.Room.UUID)
				go c.subscribeToRedisChannel(roomSubscriber, c.handleIncomingTextMessageFromRedis)

				channel = &types.Channel{
					Subscriber: roomSubscriber,
					UUID:       openRoomEvent.Room.UUID,
				}
				c.ActiveChannels[openRoomEvent.Room.UUID] = channel
			}
			channel.ParticipantsOnServer[openRoomEvent.FromUUID] = true
			c.Connections[connUUID].Conn.WriteJSON(openRoomEvent)
		}

	}

	return nil
}

func (c *MessageController) subscribeToRedisChannel(subscriber *redis.PubSub, fn func(string) error) {
	for redisMsg := range subscriber.Channel() {
		bytes, err := json.Marshal(redisMsg.Payload)
		if err != nil {
			panic(err)
		}
		err = fn(string(bytes))
		if err != nil {
			panic(err)
		}
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, contentType, Content-Type, Accept, Authorization")
}

func main() {
	r := mux.NewRouter()

	redisClient := redisClient.New()
	serverEventsSubscriber := redisClient.SetupChannel(CHANNEL_SERVER_EVENTS)

	connections := map[string]*types.Connection{}
	userConnections := map[string]map[string]bool{}
	activeChannels := map[string]*types.Channel{}

	msgController := MessageController{
		RedisClient:     &redisClient,
		Connections:     connections,
		UserConnections: userConnections,
		ActiveChannels:  activeChannels,
	}

	// subscribe to server events
	go msgController.subscribeToRedisChannel(serverEventsSubscriber, msgController.handleIncomingServerEventFromRedis)

	r.HandleFunc("/create-room", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		req := types.OpenRoomRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}

		// save the room
		roomUUID := uuid.New().String()
		room := types.ChatRoom{
			UUID:         roomUUID,
			Participants: []string{req.ToUUID, req.FromUUID},
		}

		// push this out to the redis server events channel
		openRoomEvent := types.OpenRoomEvent{
			FromUUID:  req.FromUUID,
			ToUUID:    req.ToUUID,
			EventType: EVENT_OPEN_ROOM.String(),
			Room:      room,
		}

		msgBytes, err := json.Marshal(openRoomEvent)
		if err != nil {
			panic(err)
		}

		msgController.RedisClient.PublishToRedisChannel(CHANNEL_SERVER_EVENTS, msgBytes)

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

		go msgController.setupClientConnection(conn)
	})

	fmt.Println("Opening server")
	http.ListenAndServe(":9090", r)
}

func (c *MessageController) setupClientConnection(conn *websocket.Conn) {

	var userUUID string
	var connectionUUID string
	defer func() {
		conn.Close()
		delete(c.UserConnections[userUUID], connectionUUID)
		if len(c.UserConnections) == 0 {
			delete(c.UserConnections, userUUID)
		}

		for roomUUID, channel := range c.ActiveChannels {
			_, ok := channel.ParticipantsOnServer[userUUID]
			if !ok {
				continue
			}

			// delete this client from the participants of this room
			delete(channel.ParticipantsOnServer, userUUID)

			// if no one is left on this channel, unsubscribe from it
			if len(channel.ParticipantsOnServer) == 0 {
				err := c.ActiveChannels[roomUUID].Subscriber.Close()
				if err != nil {
					panic(err)
				}
				delete(c.ActiveChannels, roomUUID)
			}
		}
	}()

	for {
		fmt.Println("GOT CONNECTION")
		// read in a message
		_, p, err := conn.ReadMessage()

		fmt.Println("GOT MESSAGE")
		fmt.Println(string(p))
		// check to see if connection has been closed
		if err != nil && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
				break
			}
		}

		msgType, err := getEventType(string(p))
		if err != nil {
			panic(err)
		}

		if msgType == EVENT_SET_CLIENT_SOCKET.String() {
			// set up the client here and send back a message to the client that everything is ready to go
			// client should be in a loading state until that happens
			msg := types.SetClientConnectionEvent{}
			err := json.Unmarshal(p, &msg)
			if err != nil {
				panic(err)
			}

			userUUID = msg.FromUUID
			connectionUUID = uuid.NewString()

			connection := &types.Connection{
				Conn: conn,
				UUID: connectionUUID,
			}

			// map the client uuid to a map of connection UUID's to the connection
			c.UserConnections[msg.FromUUID][connectionUUID] = true
			c.Connections[connectionUUID] = connection

			msg.ConnectionUUID = connectionUUID

			// send back to client the connection uuid so they can set it
			err = conn.WriteJSON(msg)
			if err != nil {
				panic(err)
			}

			// TODO - fetch existing rooms from database and subscribe
		}

		// client has sent out a text message
		if msgType == EVENT_CHAT_TEXT.String() {
			msg := types.ChatMessageEvent{}
			err := json.Unmarshal(p, &msg)
			if err != nil {
				panic(err)
			}

			// save message to database
			// push to redis channel
			c.RedisClient.PublishToRedisChannel(msg.RoomUUID, p)
		}
	}
	fmt.Println("CLOSING WEBSOCKET")
}
