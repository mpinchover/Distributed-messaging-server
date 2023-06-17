package handlers

import (
	"encoding/json"
	"messaging-service/controllers/controltower"
	"messaging-service/types/events"
	"messaging-service/types/eventtypes"
	"messaging-service/types/records"
	"messaging-service/types/requests"

	"github.com/google/uuid"
)

type Handler struct {
	ControlTowerCtrlr *controltower.ControlTowerController
}

func New() *Handler {
	controlTower := controltower.New()
	return &Handler{
		ControlTowerCtrlr: controlTower,
	}
}

func (h *Handler) getRoomsByUserUUID(req *requests.GetRoomsByUserUUIDRequest) (*requests.GetRoomsByUserUUIDResponse, error) {
	rooms, err := h.ControlTowerCtrlr.GetRoomsByUserUUID(req.UserUUID, req.Offset)
	if err != nil {
		panic(err)
	}

	h.ControlTowerCtrlr.SubscribeRoomsToServer(rooms, req.UserUUID)

	response := &requests.GetRoomsByUserUUIDResponse{
		Rooms: rooms,
	}
	return response, nil
}

func (h *Handler) getMessagesByRoomUUID(req *requests.GetMessagesByRoomUUIDRequest) (*requests.GetMessagesByRoomUUIDResponse, error) {
	msgs, err := h.ControlTowerCtrlr.GetMessagesByRoomUUID(req.RoomUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	resp := &requests.GetMessagesByRoomUUIDResponse{
		Messages: msgs,
	}
	return resp, nil
}

func (h *Handler) deleteRoom(req *requests.DeleteRoomRequest) (*requests.DeleteRoomResponse, error) {
	roomUUID := req.RoomUUID
	err := h.ControlTowerCtrlr.Repo.DeleteRoom(roomUUID)
	if err != nil {
		return nil, err
	}

	deleteRoomEvent := events.DeleteChatRoomEvent{
		EventType: eventtypes.EVENT_DELETE_ROOM.String(),
		RoomUUID:  roomUUID,
	}
	msgBytes, err := json.Marshal(deleteRoomEvent)
	if err != nil {
		return nil, err
	}

	h.ControlTowerCtrlr.RedisClient.PublishToRedisChannel(eventtypes.CHANNEL_SERVER_EVENTS, msgBytes)
	return &requests.DeleteRoomResponse{}, nil
}

func (h *Handler) createRoom(req *requests.CreateRoomRequest) (*requests.CreateRoomResponse, error) {
	// TODO, extend the 'to' field to be an array

	roomUUID := uuid.New().String()

	room := &records.ChatRoom{
		UUID: roomUUID,
		Participants: []*records.ChatParticipant{
			{
				UUID:     uuid.New().String(),
				RoomUUID: roomUUID,
				UserUUID: req.FromUUID,
			},
			{
				UUID:     uuid.New().String(),
				RoomUUID: roomUUID,
				UserUUID: req.ToUUID,
			},
		},
	}

	// push this out to the redis server events channel
	openRoomEvent := &events.OpenRoomEvent{
		FromUUID:  req.FromUUID,
		ToUUID:    req.ToUUID,
		EventType: eventtypes.EVENT_OPEN_ROOM.String(),
		Room:      room,
	}

	err := h.ControlTowerCtrlr.Repo.SaveRoom(room)
	if err != nil {
		return nil, err
	}
	msgBytes, err := json.Marshal(openRoomEvent)
	if err != nil {
		return nil, err
	}

	h.ControlTowerCtrlr.RedisClient.PublishToRedisChannel(eventtypes.CHANNEL_SERVER_EVENTS, msgBytes)
	return &requests.CreateRoomResponse{
		Room: room,
	}, nil
}
