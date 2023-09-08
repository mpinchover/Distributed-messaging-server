package handlers

import (
	"encoding/json"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"sync"
)

func (h *Handler) SetupChannels() {
	subscriber := utils.SetupChannel(h.RedisClient, enums.CHANNEL_SERVER_EVENTS)
	go utils.SubscribeToChannel(subscriber, h.HandleServerEvent)
}

func (h *Handler) HandleServerEvent(event string) error {
	eventType, err := utils.GetEventType(event)
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
	eventType, err := utils.GetEventType(event)
	if err != nil {
		return err
	}

	if eventType == enums.EVENT_TEXT_MESSAGE.String() {
		textMessageEvent := &requests.TextMessageEvent{}
		err = json.Unmarshal([]byte(event), textMessageEvent)
		if err != nil {
			return err
		}
		return h.BroadcastEventToChannelSubscribersDeviceExclusive(
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

// func (h *Handler) handleTextMessageEvent(event *requests.TextMessageEvent) error {
// 	// get the room from the server
// 	channel := h.ControlTowerCtrlr.Channels[event.Message.RoomUUID]
// 	from := event.ConnectionUUID
// 	// save the txt msg to db

// 	// room not on server
// 	if channel == nil {
// 		return errors.New("room not found")
// 	}

// 	for userUUID := range channel.MembersOnServer {
// 		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(userUUID)
// 		// issue is that connection is null
// 		for connUUID, conn := range connection.Connections {
// 			if connUUID == from {
// 				continue
// 			}
// 			conn.WriteJSON(event)
// 		}
// 	}

// 	return nil
// }

func (h *Handler) handleLeaveRoomEvent(event *requests.LeaveRoomEvent) error {
	// get the room from the server
	channel, ok := h.ControlTowerCtrlr.Channels[event.RoomUUID]
	// room not on server
	if !ok {
		return nil
	}

	// remove the user from this room
	err := h.ControlTowerCtrlr.RemoveUserFromChannel(event.UserUUID, event.RoomUUID)
	// err := h.ControlTowerCtrlr.ChannelsCtrlr.DeleteUser(event.RoomUUID, event.UserUUID)
	if err != nil {
		return err
	}

	// notify any remaining members that the user has left
	for userUUID := range channel {
		userConn := h.ControlTowerCtrlr.UserConnections[userUUID]
		for _, device := range userConn.Devices {
			device.WS.WriteJSON(event)
		}
	}
	return nil
}

func (h *Handler) handleDeleteRoomEvent(event *requests.DeleteRoomEvent) error {
	// get the room from the server
	channel, ok := h.ControlTowerCtrlr.Channels[event.RoomUUID]

	// room not on server
	if !ok {
		return nil
	}

	delete(h.ControlTowerCtrlr.Channels, event.RoomUUID)

	for userUUID := range channel {
		userConn, ok := h.ControlTowerCtrlr.UserConnections[userUUID]
		if !ok {
			continue
		}

		// notify everyone that the channel has closed
		for _, device := range userConn.Devices {
			device.WS.WriteJSON(event)
		}
	}

	return nil
}

func (h *Handler) handleOpenRoomEvent(event *requests.OpenRoomEvent) error {

	// for every member, check if they are on this server
	// if they are, then you need to subscribe the room here
	members := event.Room.Members
	roomUUID := event.Room.UUID

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// subscribe server to the room if members on are on this server
	for _, member := range members {
		userConn, ok := h.ControlTowerCtrlr.UserConnections[member.UserUUID]
		// user not on this server, so don't subscribe the room to this server
		if !ok {
			continue
		}
		// check if the room has already been subscribed to this server as well
		_, ok = h.ControlTowerCtrlr.Channels[roomUUID]
		// server contains a user who doesn't have the room subscribed
		if !ok {
			// Set up the room on this server
			subscriber := utils.SetupChannel(h.RedisClient, roomUUID)
			go utils.SubscribeToChannel(subscriber, h.HandleRoomEvent)
			h.ControlTowerCtrlr.Channels[roomUUID] = map[string]bool{}
		}

		h.ControlTowerCtrlr.Channels[roomUUID] = map[string]bool{
			member.UUID: true,
		}

		for _, device := range userConn.Devices {
			device.WS.WriteJSON(event)
		}
	}

	return nil
}
