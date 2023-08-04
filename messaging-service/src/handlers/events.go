package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"sync"

	"github.com/gorilla/websocket"
)

func (h *Handler) SetupChannels() {
	subscriber := utils.SetupChannel(h.RedisClient, enums.CHANNEL_SERVER_EVENTS)
	go utils.SubscribeToChannel(subscriber, h.HandleServerEvent)
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

func (h *Handler) HandleServerEvent(event string) error {
	eventType, err := getEventType(event)
	if err != nil {
		return err
	}

	if eventType == enums.EVENT_OPEN_ROOM.String() {
		openRoomEvent := &requests.OpenRoomEvent{}
		err = json.Unmarshal([]byte(event), openRoomEvent)
		if err != nil {
			return err
		}
		err = h.handleOpenRoomEvent(openRoomEvent)
		if err != nil {
			return err
		}

	}

	return nil
}

func (h *Handler) HandleRoomEvent(event string) error {
	eventType, err := getEventType(event)
	if err != nil {
		return err
	}

	if eventType == enums.EVENT_TEXT_MESSAGE.String() {
		textMessageEvent := &requests.TextMessageEvent{}
		err = json.Unmarshal([]byte(event), textMessageEvent)
		if err != nil {
			return err
		}
		return h.BroadcastEventToChannelSubscribersClientExclusive(
			textMessageEvent.Message.RoomUUID,
			textMessageEvent.ConnectionUUID,
			textMessageEvent,
		)
	}

	if eventType == enums.EVENT_LEAVE_ROOM.String() {
		leaveRoomEvent := &requests.LeaveRoomEvent{}
		err = json.Unmarshal([]byte(event), leaveRoomEvent)
		if err != nil {
			return err
		}

		return h.handleLeaveRoomEvent(leaveRoomEvent)
	}

	if eventType == enums.EVENT_DELETE_ROOM.String() {
		deleteRoomEvent := &requests.DeleteRoomEvent{}
		err = json.Unmarshal([]byte(event), deleteRoomEvent)
		if err != nil {
			return err
		}
		// return h.BroadcastEventToChannelSubscribersExclusive(
		// 	deleteMessageEvent.RoomUUID,
		// 	deleteMessageEvent.UserUUID,
		// 	deleteMessageEvent,
		// )
		return h.handleDeleteRoomEvent(deleteRoomEvent)
	}

	if eventType == enums.EVENT_SEEN_MESSAGE.String() {
		seenMsgEvent := &requests.SeenMessageEvent{}
		err = json.Unmarshal([]byte(event), seenMsgEvent)
		if err != nil {
			return err
		}
		return h.BroadcastEventToChannelSubscribersUserExclusive(seenMsgEvent.RoomUUID, seenMsgEvent.UserUUID, seenMsgEvent)
	}

	if eventType == enums.EVENT_DELETE_MESSAGE.String() {
		deleteMessageEvent := &requests.DeleteMessageEvent{}
		err = json.Unmarshal([]byte(event), deleteMessageEvent)
		if err != nil {
			return err
		}
		return h.BroadcastEventToChannelSubscribers(
			deleteMessageEvent.RoomUUID,
			deleteMessageEvent,
		)
	}

	return nil
}

func (h *Handler) handleTextMessageEvent(event *requests.TextMessageEvent) error {
	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(event.Message.RoomUUID)
	from := event.ConnectionUUID
	// save the txt msg to db

	// room not on server
	if channel == nil {
		return errors.New("room not found")
	}

	for userUUID := range channel.MembersOnServer {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(userUUID)
		// issue is that connection is null
		for connUUID, conn := range connection.Connections {
			if connUUID == from {
				continue
			}
			conn.WriteJSON(event)
		}
	}

	return nil
}

func (h *Handler) handleLeaveRoomEvent(event *requests.LeaveRoomEvent) error {
	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(event.RoomUUID)

	// room not on server
	if channel == nil {
		return nil
	}

	// remove the user from this room
	err := h.ControlTowerCtrlr.ChannelsCtrlr.DeleteUser(event.RoomUUID, event.UserUUID)
	if err != nil {
		return err
	}

	// notify any remaining members that the user has left
	for userUUID := range channel.MembersOnServer {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(userUUID)
		for _, conn := range connection.Connections {
			conn.WriteJSON(event)
		}
	}
	return nil
}

func (h *Handler) handleDeleteRoomEvent(event *requests.DeleteRoomEvent) error {
	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(event.RoomUUID)

	// room not on server
	if channel == nil {
		return nil
	}

	h.ControlTowerCtrlr.ChannelsCtrlr.DeleteChannel(channel.UUID)

	// notify everyone that the channel has closed
	for userUUID := range channel.MembersOnServer {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(userUUID)
		for _, conn := range connection.Connections {
			conn.WriteJSON(event)
		}
	}

	return nil
}

func (h *Handler) handleOpenRoomEvent(event *requests.OpenRoomEvent) error {

	// for every member, check if they are on this server
	// if they are, then you need to subscribe the room here
	members := event.Room.Members
	roomUUID := event.Room.UUID

	connectionsToWrite := []*websocket.Conn{}
	for _, member := range members {

		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(member.UserUUID)
		// user not on this server, so don't subscribe the room to this server
		if connection == nil {
			continue
		}
		// check if the room has already been subscribed to this server as well
		room := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(roomUUID)

		if room == nil {
			// Set up the room on this server
			subscriber := utils.SetupChannel(h.RedisClient, roomUUID)
			go utils.SubscribeToChannel(subscriber, h.HandleRoomEvent)

			room = &requests.ServerChannel{
				Subscriber:      subscriber,
				MembersOnServer: map[string]bool{},
				UUID:            roomUUID,
			}
			h.ControlTowerCtrlr.ChannelsCtrlr.AddChannel(room)
		}

		h.ControlTowerCtrlr.ChannelsCtrlr.AddUserToChannel(room.UUID, member.UserUUID)

		// for all connections, write the openRoomEvent
		for _, conn := range connection.Connections {
			connectionsToWrite = append(connectionsToWrite, conn)
		}
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range connectionsToWrite {
		err := conn.WriteJSON(event)
		if err != nil {
			log.Println("Failed to write message back to connection")
		}
	}

	return nil
}
