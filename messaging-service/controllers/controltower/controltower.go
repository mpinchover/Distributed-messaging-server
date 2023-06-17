package controltower

import (
	"encoding/json"
	"errors"
	"fmt"
	redisClient "messaging-service/redis"
	"messaging-service/repo"
	"messaging-service/types/entities"
	"messaging-service/types/events"
	"messaging-service/types/eventtypes"
	"messaging-service/types/records"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type ControlTowerController struct {
	RedisClient *redisClient.RedisClient
	Connections map[string]*entities.Connection

	MapLock *sync.Mutex

	// map the user uuid to a list of the user's connections (different devices)
	UserConnections         map[string]map[string]bool
	OutboundMessagesChannel <-chan *redis.Message

	// track active rooms/channels on this server
	ActiveChannels map[string]*entities.Channel
	Repo           *repo.Repo
}

func New() *ControlTowerController {
	redisClient := redisClient.New()

	repo, err := repo.New()
	if err != nil {
		panic(err)
	}

	connections := map[string]*entities.Connection{}
	userConnections := map[string]map[string]bool{}
	activeChannels := map[string]*entities.Channel{}

	var mu sync.Mutex
	msgController := &ControlTowerController{
		RedisClient:     &redisClient,
		Connections:     connections,
		UserConnections: userConnections,
		ActiveChannels:  activeChannels,
		Repo:            repo,
		MapLock:         &mu,
	}

	serverEventsSubscriber := redisClient.SetupChannel(eventtypes.CHANNEL_SERVER_EVENTS)
	go msgController.subscribeToRedisChannel(serverEventsSubscriber, msgController.handleIncomingServerEventFromRedis)

	return msgController
}

func (c *ControlTowerController) GetMessagesByRoomUUID(roomUUID string, offset int) ([]*records.ChatMessage, error) {
	return c.Repo.GetMessagesByRoomUUID(roomUUID, offset)
}

func (c *ControlTowerController) PublishAndSubscribeRoom(userUUID string, roomUUID string, openRoomEvent *events.OpenRoomEvent) error {
	listOfFromConnections, ok := c.UserConnections[userUUID]

	// no connections for this user uuid found on this server
	if !ok {
		return nil
	}

	for connUUID := range listOfFromConnections {
		c.MapLock.Lock()
		channel, ok := c.ActiveChannels[roomUUID]
		if !ok {
			roomSubscriber := c.RedisClient.SetupChannel(roomUUID)
			go c.subscribeToRedisChannel(roomSubscriber, c.handleIncomingTextMessageFromRedis)

			channel = &entities.Channel{
				Subscriber:           roomSubscriber,
				UUID:                 roomUUID,
				ParticipantsOnServer: map[string]bool{},
			}
			c.ActiveChannels[roomUUID] = channel
		}
		channel.ParticipantsOnServer[userUUID] = true

		c.Connections[connUUID].Conn.WriteJSON(openRoomEvent)
		c.MapLock.Unlock()
	}
	return nil
}

func (c *ControlTowerController) handleIncomingServerEventFromRedis(event string) error {
	eventType, err := getEventType(event)
	if err != nil {
		panic(err)
	}

	if eventType == eventtypes.EVENT_OPEN_ROOM.String() {

		openRoomEvent := &events.OpenRoomEvent{}
		err = json.Unmarshal([]byte(event), openRoomEvent)
		if err != nil {
			panic(err)
		}

		fromUUID := openRoomEvent.FromUUID
		toUUID := openRoomEvent.ToUUID
		roomUUID := openRoomEvent.Room.UUID

		c.PublishAndSubscribeRoom(fromUUID, roomUUID, openRoomEvent)
		c.PublishAndSubscribeRoom(toUUID, roomUUID, openRoomEvent)
	}
	return nil
}

func (c *ControlTowerController) subscribeToRedisChannel(subscriber *redis.PubSub, fn func(string) error) {
	for redisMsg := range subscriber.Channel() {
		err := fn(redisMsg.Payload)
		if err != nil {
			panic(err)
		}
	}
}

func (c *ControlTowerController) handleIncomingTextMessageFromRedis(msg string) error {
	chatMessage := events.ChatMessageEvent{}
	err := json.Unmarshal([]byte(msg), &chatMessage)
	if err != nil {
		panic(err)
	}

	roomUUID := chatMessage.RoomUUID
	room, ok := c.ActiveChannels[roomUUID]
	if !ok {
		return nil
	}

	// get all the outbound connections we need to send the message
	outboundConnections := []*entities.Connection{}
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

	c.MapLock.Lock()
	for _, outboundConn := range outboundConnections {
		outboundConn.Conn.WriteJSON(chatMessage)
	}
	c.MapLock.Unlock()

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

func (c *ControlTowerController) SetupClientConnection(conn *websocket.Conn) {

	var userUUID string
	var connectionUUID string
	conn.SetPongHandler(func(appData string) error {
		err := conn.WriteMessage(1, []byte("PONG"))
		if err != nil {
			panic(err)
		}
		return nil
	})

	defer func() {
		conn.Close()
		delete(c.UserConnections[userUUID], connectionUUID)
		if len(c.UserConnections) == 0 {
			delete(c.UserConnections, userUUID)
		}

		// log.Println("HI THERE!!")
		c.MapLock.Lock()
		defer func() {
			c.MapLock.Unlock()
		}()
		// TODO - move this to a channel
		for roomUUID, channel := range c.ActiveChannels {
			_, ok := channel.ParticipantsOnServer[userUUID]
			if !ok {
				continue
			}

			// if two clients are attached to the same server, they will both try to delete from
			// the same map

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
		// read in a message
		_, p, err := conn.ReadMessage()

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

		if msgType == eventtypes.EVENT_SET_CLIENT_SOCKET.String() {
			// set up the client here and send back a message to the client that everything is ready to go
			// client should be in a loading state until that happens

			msg := events.SetClientConnectionEvent{}
			err := json.Unmarshal(p, &msg)
			if err != nil {
				panic(err)
			}

			userUUID = msg.FromUUID
			connectionUUID := uuid.New().String()

			connection := &entities.Connection{
				Conn: conn,
				UUID: connectionUUID,
			}

			// map the client uuid to a map of connection UUID's to the connection
			_, ok := c.UserConnections[userUUID]
			if !ok {
				c.UserConnections[userUUID] = map[string]bool{}
			}
			c.UserConnections[userUUID][connectionUUID] = true

			c.Connections[connectionUUID] = connection

			msg.ConnectionUUID = connectionUUID

			// send back to client the connection uuid so they can set it
			err = conn.WriteJSON(msg)
			if err != nil {
				panic(err)
			}

		}

		// client has sent out a text message
		if msgType == eventtypes.EVENT_CHAT_TEXT.String() {

			go c.handleTextMessage(p)
		}
	}
}

func (c *ControlTowerController) handleTextMessage(p []byte) {
	msg := events.ChatMessageEvent{}
	err := json.Unmarshal(p, &msg)
	if err != nil {
		panic(err)
	}

	chatMessage := &records.ChatMessage{
		FromUUID:    msg.FromUserUUID,
		MessageText: msg.MessageText,
		RoomUUID:    msg.RoomUUID,
		UUID:        uuid.New().String(),
	}

	err = c.Repo.SaveChatMessage(chatMessage)
	if err != nil {
		panic(err)
	}

	roomUUID := msg.RoomUUID
	c.RedisClient.PublishToRedisChannel(roomUUID, p)
}

func (c *ControlTowerController) GetRoomsByUserUUID(userUUID string, offset int) ([]*records.ChatRoom, error) {
	return c.Repo.GetRoomsByUserUUID(userUUID, offset)
}

func (c *ControlTowerController) SubscribeRoomsToServer(rooms []*records.ChatRoom, userUUID string) {
	for _, room := range rooms {
		roomUUID := room.UUID
		_, ok := c.ActiveChannels[roomUUID]
		if ok {
			continue
		}
		// if we are not already subscribed to the channel on this server, do so.

		roomSubscriber := c.RedisClient.SetupChannel(roomUUID)
		go c.subscribeToRedisChannel(roomSubscriber, c.handleIncomingTextMessageFromRedis)

		channel := &entities.Channel{
			Subscriber:           roomSubscriber,
			UUID:                 roomUUID,
			ParticipantsOnServer: map[string]bool{},
		}
		c.ActiveChannels[roomUUID] = channel
		channel.ParticipantsOnServer[userUUID] = true
	}
}
