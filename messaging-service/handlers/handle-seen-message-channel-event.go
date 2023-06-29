package handlers

import (
	"log"
	"messaging-service/types/requests"
)

func (h *Handler) handleSeenMessageChannelEvent(msg *requests.SeenMessageEvent) error {
	// if the user connection is on this server, blast it out.
	members, err := h.ControlTowerCtrlr.Repo.GetMembersByRoomUUID(msg.RoomUUID)
	if err != nil {
		return err
	}

	for _, member := range members {
		connection := h.ControlTowerCtrlr.ConnCtrlr.GetConnection(member.UserUUID)
		if connection == nil {
			continue
		}

		for _, conn := range connection.Connections {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
