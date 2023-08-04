package handlers

import (
	"errors"
)

// broadcast to all channel members excluding the client device
func (h *Handler) BroadcastEventToChannelSubscribersClientExclusive(channelUUID string, fromClientUUID string, msg interface{}) error {

	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(channelUUID)

	// room not on server
	if channel == nil {
		return errors.New("room not found")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	for _, m := range members {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(m.UserUUID)
		if connection == nil {
			continue
		}
		for clientUUID, conn := range connection.Connections {
			if clientUUID == fromClientUUID {
				continue
			}
			conn.WriteJSON(msg)
		}
	}
	return nil
}

// broadcast to all channel members excluding any userUUID devices
func (h *Handler) BroadcastEventToChannelSubscribersUserExclusive(channelUUID string, userUUID string, msg interface{}) error {

	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(channelUUID)

	// room not on server
	if channel == nil {
		return errors.New("room not found")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	for _, m := range members {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(m.UserUUID)

		// don't broadcast to any devices belonging to this user
		if m.UserUUID == userUUID {
			continue
		}

		for _, conn := range connection.Connections {
			conn.WriteJSON(msg)
		}
	}
	return nil
}

// broadcast to all channel members excluding any userUUID devices
func (h *Handler) BroadcastEventToChannelSubscribers(channelUUID string, msg interface{}) error {

	// get the room from the server
	channel := h.ControlTowerCtrlr.ChannelsCtrlr.GetChannel(channelUUID)

	// room not on server
	if channel == nil {
		return errors.New("room not found")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	for _, m := range members {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(m.UserUUID)

		for _, conn := range connection.Connections {
			conn.WriteJSON(msg)
		}
	}
	return nil
}
