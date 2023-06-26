package controltower

import (
	"encoding/json"
	"errors"
	redisClient "messaging-service/redis"
	"messaging-service/repo"
	"messaging-service/types/enums"
	"messaging-service/types/records"
	"messaging-service/types/requests"
	"messaging-service/utils"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ControlTowerController struct {
	RedisClient *redisClient.RedisClient
	Connections map[string]*requests.Connection
	// track active rooms/channels on this server
	ServerChannels map[string]*requests.ServerChannel

	MapLock *sync.Mutex

	Repo *repo.Repo
}

func New(
	redisClient *redisClient.RedisClient,
	repo *repo.Repo,
) *ControlTowerController {
	connections := map[string]*requests.Connection{}
	serverChannels := map[string]*requests.ServerChannel{}

	var mu sync.Mutex
	msgController := &ControlTowerController{
		RedisClient:    redisClient,
		Connections:    connections,
		ServerChannels: serverChannels,

		Repo:    repo,
		MapLock: &mu,
	}

	return msgController
}

func (c *ControlTowerController) GetMessagesByRoomUUID(roomUUID string, offset int) ([]*records.Message, error) {
	return c.Repo.GetMessagesByRoomUUID(roomUUID, offset)
}

func (c *ControlTowerController) CreateRoom(
	members []*requests.Member,
) (*requests.Room, error) {
	// build the room
	for _, member := range members {
		member.UUID = uuid.New().String()
	}

	roomUUID := uuid.New().String()
	repoMembers := make([]*records.Member, len(members))

	// TODO – if member is nil, make it a default type
	for i, member := range members {
		repoMembers[i] = &records.Member{
			UUID:     member.UUID,
			RoomUUID: roomUUID,
			UserUUID: member.UserUUID,
			UserRole: member.UserRole,
		}
	}

	repoRoom := &records.Room{
		UUID:    roomUUID,
		Members: repoMembers,
	}

	err := c.Repo.SaveRoom(repoRoom)
	if err != nil {
		return nil, err
	}

	newRoom := &requests.Room{
		Members: members,
		UUID:    roomUUID,
	}

	openRoomEvent := requests.OpenRoomEvent{
		EventType: enums.EVENT_OPEN_ROOM.String(),
		Room:      newRoom,
	}

	err = utils.PublishToRedisChannel(c.RedisClient, enums.CHANNEL_SERVER_EVENTS, openRoomEvent)
	if err != nil {
		return nil, err
	}
	return newRoom, nil
}

func (c *ControlTowerController) LeaveRoom(userUUID string, roomUUID string) error {
	room, err := c.Repo.GetRoomByRoomUUID(roomUUID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room does not exist")
	}

	// TODO – this is something the client should verify not the server
	// membersInRoom := make([]string, len(room.Members))
	// for i, mem := range room.Members {
	// 	membersInRoom[i] = mem.UserUUID
	// }

	// if !utils.Contains(membersInRoom, userUUID) {
	// 	return errors.New("member not in room")
	// }

	// TODO - put in helper function
	// TODO – in the future add in fn to make this optional
	if len(room.Members) == 1 {
		err := c.Repo.DeleteRoom(roomUUID)
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
		c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
		return nil
	}

	err = c.Repo.LeaveRoom(userUUID, roomUUID)
	if err != nil {
		return err
	}
	leaveRoomEvent := requests.LeaveRoomEvent{
		EventType: enums.EVENT_LEAVE_ROOM.String(),
		RoomUUID:  roomUUID,
		UserUUID:  userUUID,
	}
	msgBytes, err := json.Marshal(leaveRoomEvent)
	if err != nil {
		return err
	}

	c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
	return nil
}

func (c *ControlTowerController) DeleteRoom(roomUUID string) error {
	room, err := c.Repo.GetRoomByRoomUUID(roomUUID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	// put in helper function
	membersInRoom := make([]string, len(room.Members))
	for i, mem := range room.Members {
		membersInRoom[i] = mem.UserUUID
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

	c.RedisClient.PublishToRedisChannel(roomUUID, msgBytes)
	return nil
}

func (c *ControlTowerController) SetupClientConnectionV2(
	conn *websocket.Conn,
	msg *requests.SetClientConnectionEvent) (*requests.SetClientConnectionEvent, error) {

	connectionUUID := uuid.New().String()
	msg.ConnectionUUID = connectionUUID
	connection := c.GetClientConnectionFromServer(msg.FromUUID)

	if connection == nil {
		connection = &requests.Connection{
			UserUUID:    msg.FromUUID,
			Connections: map[string]*websocket.Conn{},
		}
		c.AddUserConnection(connection)
	}

	c.AddClientConnection(connection, connectionUUID, conn)
	return msg, nil
}

func (c *ControlTowerController) GetRoomsByUserUUID(userUUID string, offset int) ([]*requests.Room, error) {
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
				UserRole: member.UserRole,
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
