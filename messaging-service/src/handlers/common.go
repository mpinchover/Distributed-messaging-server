package handlers

import (
	"errors"
	"sync"
)

// broadcast to all channel members excluding the client device
func (h *Handler) BroadcastEventToChannelSubscribersDeviceExclusive(channelUUID string, fromDeviceUUID string, msg interface{}) error {

	// get the room from the server
	_, ok := h.ControlTowerCtrlr.Channels[channelUUID]
	if !ok {
		return errors.New("room not found on server")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for _, m := range members {
		userConn, ok := h.ControlTowerCtrlr.UserConnections[m.UserUUID]
		if !ok {
			continue
		}
		for deviceUUID, device := range userConn.Devices {
			if deviceUUID == fromDeviceUUID {
				continue
			}

			device.WS.WriteJSON(msg)

		}
	}
	return nil
}

// broadcast to all channel members excluding any userUUID devices
func (h *Handler) BroadcastEventToChannelSubscribersUserExclusive(channelUUID string, userUUID string, msg interface{}) error {

	// get the room from the server
	_, ok := h.ControlTowerCtrlr.Channels[channelUUID]

	// room not on server
	if !ok {
		return errors.New("room not found on server")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for _, m := range members {
		userConn, ok := h.ControlTowerCtrlr.UserConnections[m.UserUUID]
		if !ok {
			continue
		}

		// don't broadcast to any devices belonging to this user
		if m.UserUUID == userUUID {
			continue
		}

		for _, device := range userConn.Devices {
			device.WS.WriteJSON(msg)
		}
	}
	return nil
}

// broadcast to all channel members excluding any userUUID devices
func (h *Handler) BroadcastEventToChannelSubscribers(channelUUID string, msg interface{}) error {

	// get the room from the server
	_, ok := h.ControlTowerCtrlr.UserConnections[channelUUID]
	// room not on server
	if !ok {
		return errors.New("room not found on server")
	}

	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(channelUUID)
	if err != nil {
		return err
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for _, m := range members {
		userConn, ok := h.ControlTowerCtrlr.UserConnections[m.UserUUID]
		if !ok {
			continue
		}

		for _, device := range userConn.Devices {
			device.WS.WriteJSON(msg)
		}
	}
	return nil
}
