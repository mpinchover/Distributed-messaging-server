package handlers

import (
	"messaging-service/controllers/controltower"
	redisClient "messaging-service/redis"
	"messaging-service/types/requests"
	"messaging-service/validation"
)

type Handler struct {
	ControlTowerCtrlr *controltower.ControlTowerController
	RedisClient       *redisClient.RedisClient
}

func New() *Handler {

	redisClient := redisClient.New()
	// subscrie to the events here for server
	controlTower := controltower.New()

	return &Handler{
		ControlTowerCtrlr: controlTower,
		RedisClient:       &redisClient,
	}
}

func (h *Handler) getRoomsByUserUUID(req *requests.GetRoomsByUserUUIDRequest) (*requests.GetRoomsByUserUUIDResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	rooms, err := h.ControlTowerCtrlr.GetRoomsByUserUUID(req.UserUUID, req.Offset)
	if err != nil {
		panic(err)
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

	response := &requests.GetRoomsByUserUUIDResponse{
		Rooms: requestRooms,
	}
	return response, nil
}

func (h *Handler) getMessagesByRoomUUID(req *requests.GetMessagesByRoomUUIDRequest) (*requests.GetMessagesByRoomUUIDResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	msgs, err := h.ControlTowerCtrlr.GetMessagesByRoomUUID(req.RoomUUID, req.Offset)
	if err != nil {
		return nil, err
	}

	requestMsgs := make([]*requests.Message, len(msgs))

	for i, msg := range msgs {
		requestMsgs[i] = &requests.Message{
			UUID:        msg.UUID,
			FromUUID:    msg.FromUUID,
			RoomUUID:    msg.RoomUUID,
			MessageText: msg.MessageText,
			CreatedAt:   msg.Model.CreatedAt.UnixMilli(),
		}
	}

	resp := &requests.GetMessagesByRoomUUIDResponse{
		Messages: requestMsgs,
	}
	return resp, nil
}

func (h *Handler) deleteRoom(req *requests.DeleteRoomRequest) (*requests.DeleteRoomResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	roomUUID := req.RoomUUID
	// verify user has permissions
	err = h.ControlTowerCtrlr.DeleteRoom(roomUUID, req.UserUUID)
	if err != nil {
		return nil, err
	}
	return &requests.DeleteRoomResponse{}, nil
}

func (h *Handler) createRoom(req *requests.CreateRoomRequest) (*requests.CreateRoomResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	room, err := h.ControlTowerCtrlr.CreateRoom(req.Members)
	if err != nil {
		return nil, err
	}

	return &requests.CreateRoomResponse{
		Room: room,
	}, nil
}

func (h *Handler) leaveRoom(req *requests.LeaveRoomRequest) (*requests.LeaveRoomResponse, error) {
	err := validation.ValidateRequest(req)
	if err != nil {
		return nil, err
	}

	err = h.ControlTowerCtrlr.LeaveRoom(req.UserUUID, req.RoomUUID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
