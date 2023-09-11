package controltower

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	redisClient "messaging-service/src/redis"
	"messaging-service/src/repo"
	"messaging-service/src/serrors"
	"messaging-service/src/types/connections"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ControlTowerCtrlr struct {
	RedisClient redisClient.RedisInterface
	Repo        repo.RepoInterface

	UserConnections map[string]*connections.UserConnection // user uuid to list of devices
	Channels        map[string]*connections.Channel        // room uuid to the list of users in the room
}

func New(
	redisClient *redisClient.RedisClient,
	repo *repo.Repo,
) *ControlTowerCtrlr {

	controlTower := &ControlTowerCtrlr{
		RedisClient:     redisClient,
		Repo:            repo,
		Channels:        map[string]*connections.Channel{},
		UserConnections: map[string]*connections.UserConnection{},
	}

	return controlTower
}

func (c *ControlTowerCtrlr) GetMessagesByRoomUUID(ctx context.Context, roomUUID string, offset int) ([]*records.Message, error) {
	return c.Repo.GetMessagesByRoomUUID(roomUUID, offset)

}


func (c *ControlTowerCtrlr) CreateRoom(
	ctx context.Context,
	members []*requests.Member,
) (*requests.Room, error) {
	for _, member := range members {
		member.UUID = uuid.New().String()
	}

	roomUUID := uuid.New().String()
	repoMembers := make([]*records.Member, len(members))

	for i, member := range members {
		repoMembers[i] = &records.Member{
			UUID:     member.UUID,
			RoomUUID: roomUUID,
			UserUUID: member.UserUUID,
		}
	}

	createdAtNano := float64(time.Now().UnixNano()) //  1e6
	repoRoom := &records.Room{
		UUID:          roomUUID,
		Members:       repoMembers,
		CreatedAtNano: createdAtNano,
	}

	// go func(room *records.Room) {
	// 	err := c.Repo.SaveRoom(repoRoom)
	// 	if err != nil {
	// 		log.Println("failed to save room")
	// 	}
	// }(repoRoom)
	err := c.Repo.SaveRoom(repoRoom)
	if err != nil {
		return nil, err
	}

	newRoom := &requests.Room{
		Members:       members,
		UUID:          roomUUID,
		CreatedAtNano: createdAtNano,
	}

	openRoomEvent := requests.OpenRoomEvent{
		EventType: enums.EVENT_OPEN_ROOM.String(),
		Room:      newRoom,
	}

	bytes, err := json.Marshal(openRoomEvent)
	if err != nil {
		return nil, err
	}

	err = c.RedisClient.PublishToRedisChannel(enums.CHANNEL_SERVER_EVENTS, bytes)
	if err != nil {
		return nil, err
	}
	return newRoom, nil
}

func (c *ControlTowerCtrlr) UpdateMessage(ctx context.Context, message *requests.Message) error {
	// first get the message
	existingMsg, err := c.Repo.GetMessageByUUID(message.UUID)
	if err != nil {
		return err
	}

	// if we haven't already deleted the message and want to delete it
	if existingMsg.MessageStatus != enums.MESSAGE_STATUS_DELETED.String() &&
		message.MessageStatus != existingMsg.MessageStatus {
		existingMsg.MessageStatus = message.MessageStatus
	}

	return c.Repo.UpdateMessage(existingMsg)
}

// func (c *ControlTowerCtrlr) LeaveRoom(ctx context.Context, userUUID string, roomUUID string) error {
// 	room, err := c.Repo.GetRoomByRoomUUID(roomUUID)
// 	if err != nil {
// 		return err
// 	}
// 	if room == nil {
// 		return errors.New("room does not exist")
// 	}

// 	// TODO – this is something the client should verify not the server

// 	// TODO - put in helper function
// 	// TODO – in the future add in fn to make this optional
// 	if len(room.Members) == 1 {
// 		err := c.Repo.DeleteRoom(roomUUID)
// 		if err != nil {
// 			return err
// 		}
// 		deleteRoomEvent := requests.DeleteRoomEvent{
// 			EventType: enums.EVENT_DELETE_ROOM.String(),
// 			RoomUUID:  roomUUID,
// 		}

// 		msgBytes, err := json.Marshal(deleteRoomEvent)
// 		if err != nil {
// 			return err
// 		}
// 		return c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
// 	}

// 	err = c.Repo.LeaveRoom(userUUID, roomUUID)
// 	if err != nil {
// 		return err
// 	}
// 	leaveRoomEvent := requests.LeaveRoomEvent{
// 		EventType: enums.EVENT_LEAVE_ROOM.String(),
// 		RoomUUID:  roomUUID,
// 		UserUUID:  userUUID,
// 	}
// 	msgBytes, err := json.Marshal(leaveRoomEvent)
// 	if err != nil {
// 		return err
// 	}

// 	return c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
// }

func (c *ControlTowerCtrlr) DeleteRoom(ctx context.Context, roomUUID string) error {
	room, err := c.Repo.GetRoomByRoomUUID(roomUUID)
	if err != nil {
		return err
	}
	if room == nil {
		return serrors.InternalErrorf("room not found", nil)
	}

	err = c.Repo.DeleteRoom(roomUUID)
	if err != nil {
		return err
	}

	deleteRoomEvent := requests.DeleteRoomEvent{
		EventType: enums.EVENT_DELETE_ROOM.String(),
		RoomUUID:  roomUUID,
	}
	msgBytes, err := json.Marshal(deleteRoomEvent)
	if err != nil {
		return err
	}

	return c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
}

func (c *ControlTowerCtrlr) SetupClientConnectionV2(
	conn *websocket.Conn,
	msg *requests.SetClientConnectionEvent) (*requests.SetClientConnectionEvent, error) {

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	deviceUUID := uuid.New().String()
	msg.DeviceUUID = deviceUUID
	userConn, ok := c.UserConnections[msg.UserUUID]
	if !ok {
		userConn = &connections.UserConnection{
			UUID:    msg.UserUUID,
			Devices: map[string]*connections.Device{},
		}
		c.UserConnections[msg.UserUUID] = userConn
	}
	newDeviceConnection := &connections.Device{
		WS: conn,
	}
	userConn.Devices[deviceUUID] = newDeviceConnection
	c.UserConnections[msg.UserUUID] = userConn

	return msg, nil
}

func (c *ControlTowerCtrlr) GetRoomsByUserUUIDForSubscribing(userUUID string) ([]*records.Room, error) {
	rooms, err := c.Repo.GetRoomsByUserUUIDForSubscribing(userUUID)
	return rooms, err
}

func (c *ControlTowerCtrlr) SaveSeenBy(msg *requests.SeenMessageEvent) error {

	// todo, put in concurrent fn
	existingMessage, err := c.Repo.GetMessageByUUID(msg.MessageUUID)
	if err != nil {
		return err
	}

	if existingMessage == nil {
		return errors.New("message not found")
	}

	seenBy := &records.SeenBy{
		UserUUID:    msg.UserUUID,
		MessageID:   int(existingMessage.Model.ID),
		MessageUUID: msg.MessageUUID,
	}

	err = c.Repo.SaveSeenBy(seenBy)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.RedisClient.PublishToRedisChannel(msg.RoomUUID, bytes)
}

func (c *ControlTowerCtrlr) GetRoomsByUserUUID(ctx context.Context, userUUID string, offset int) ([]*requests.Room, error) {
	rooms, err := c.Repo.GetRoomsByUserUUID(userUUID, offset)
	if err != nil {
		return nil, err
	}

	// TODO - put this all in the controller
	requestRooms := make([]*requests.Room, len(rooms))
	for i, room := range rooms {
		members := make([]*requests.Member, len(room.Members))
		messages := make([]*requests.Message, len(room.Messages))

		for j, member := range room.Members {
			members[j] = &requests.Member{
				UserUUID: member.UserUUID,
			}
		}

		for j, msg := range room.Messages {
			messages[j] = &requests.Message{
				UUID:        msg.UUID,
				FromUUID:    msg.FromUUID,
				RoomUUID:    msg.RoomUUID,
				MessageText: msg.MessageText,
			}
		}

		requestRooms[i] = &requests.Room{
			UUID:     room.UUID,
			Members:  members,
			Messages: messages,
		}
	}
	return requestRooms, nil
}

// // refer to removing the client device
// // don't need this, prob just use delete room
// func (c *ControlTowerCtrlr) RemoveUserFromChannel(userUUID string, channelUUID string) error {
// 	ch, ok := c.Channels[channelUUID]
// 	if !ok {
// 		return nil
// 	}

// 	var mu sync.Mutex
// 	mu.Lock()
// 	defer mu.Unlock()

// 	delete(ch.Users, userUUID)
// 	c.Channels[channelUUID] = ch

// 	if len(ch.Users) == 1 {
// 		delete(c.Channels, channelUUID)
// 		err := ch.Subscriber.Unsubscribe(context.Background())
// 		if err != nil {
// 			// TODO - remove the error
// 			panic(err)
// 		}
// 	}
// 	return nil
// }

// maybe store the rooms each member is part of as memebersOnServer
// remove device from server
// todo, can ou just store the channel in mysql/redis?
func (c *ControlTowerCtrlr) RemoveClientDeviceFromServer(userUUID string, deviceUUID string) error {
	// remove the user from connections

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	userConnection, ok := c.UserConnections[userUUID]
	if !ok {
		panic("User not found in user connections")
	}

	_, ok = userConnection.Devices[deviceUUID]
	if !ok {
		panic("device not in device connections")
	}

	delete(userConnection.Devices, deviceUUID)

	// user has no more devices attached to this connection, delete it
	if len(userConnection.Devices) == 0 {
		delete(c.UserConnections, userUUID)
	} else {
		// otherwise reset the user connections
		c.UserConnections[userUUID] = userConnection
	}

	// get the channels for this user, possible optimization
	// rooms, err := c.GetRoomsByUserUUIDForSubscribing(userUUID)
	// if err != nil {
	// 	return err
	// }

	// for _, ch := range rooms {
	// 	_, ok := c.Channels[]
	// }

	// iterate over every channel
	for chUUID, ch := range c.Channels {
		// if the user is not in this channel, continue
		if !ch.Users[userUUID] {
			continue
		}

		// check to see if we have removed the user
		// if the user had no more devices connected, we removed them
		_, userHasDevicesConnected := c.UserConnections[userUUID]

		// user has been deleted; delete user from the channel on this server
		if !userHasDevicesConnected {
			delete(ch.Users, userUUID)
		}

		// if no one else is in channel, unsubscribe and delete channel
		if len(ch.Users) == 0 {
			delete(c.Channels, chUUID)
			err := ch.Subscriber.Unsubscribe(context.Background())
			if err != nil {
				// TODO - remove the panic
				panic(err)
			}
		} else {
			// otherwise, update the channel
			c.Channels[chUUID] = ch
		}
	}

	return nil
}

// for testing only, add an admin token
func (c *ControlTowerCtrlr) GetUserConnections() map[string]*connections.UserConnection {
	// userConn := c.UserConnections[userUUID]
	return c.UserConnections
}

func (c *ControlTowerCtrlr) GetChannel() map[string]*connections.Channel {
	return c.Channels
	// _, ok := c.Channels[chUUID]
	// if !ok {
	// 	return map[string]bool{}
	// }
	// return c.Channels[chUUID].Users
}
