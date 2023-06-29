package controltower

import (
	"encoding/json"
	"errors"
	"messaging-service/mappers"
	"messaging-service/types/records"
	"messaging-service/types/requests"

	"github.com/google/uuid"
)

func (c *ControlTowerCtrlr) ProcessTextMessage(msg *requests.TextMessageEvent) (*requests.Message, error) {
	// ensure room exists
	room, err := c.Repo.GetRoomByRoomUUID(msg.RoomUUID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room does not exist")
	}

	msgUUID := uuid.New().String()
	msg.MessageUUID = msgUUID

	repoMessage := &records.Message{
		FromUUID:    msg.FromUUID,
		RoomUUID:    msg.RoomUUID,
		RoomID:      int(room.Model.ID),
		MessageText: msg.MessageText,
		UUID:        msgUUID,
	}

	err = c.Repo.SaveMessage(repoMessage)
	if err != nil {
		return nil, err
	}

	requestsMessage := mappers.FromRecordsMessageToRequestMessage(repoMessage)
	msg.CreatedAt = requestsMessage.CreatedAt

	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	c.RedisClient.PublishToRedisChannel(msg.RoomUUID, bytes)
	return requestsMessage, nil
}
