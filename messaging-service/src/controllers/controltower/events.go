package controltower

import (
	"encoding/json"
	"errors"
	"messaging-service/src/mappers"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"

	"github.com/google/uuid"
)

// TODO – event should just have the message embedded within it
func (c *ControlTowerCtrlr) ProcessTextMessage(msg *requests.TextMessageEvent) (*requests.Message, error) {
	// ensure room exists
	room, err := c.Repo.GetRoomByRoomUUID(msg.Message.RoomUUID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room does not exist")
	}

	msgUUID := uuid.New().String()
	msg.Message.UUID = msgUUID

	repoMessage := &records.Message{
		FromUUID:      msg.FromUUID,
		RoomUUID:      msg.Message.RoomUUID,
		RoomID:        int(room.Model.ID),
		MessageText:   msg.Message.MessageText,
		UUID:          msgUUID,
		MessageStatus: enums.MESSAGE_STATUS_LIVE.String(),
	}

	err = c.Repo.SaveMessage(repoMessage)
	if err != nil {
		return nil, err
	}

	requestsMessage := mappers.FromRecordsMessageToRequestMessage(repoMessage)
	msg.Message.CreatedAt = requestsMessage.CreatedAt

	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	c.RedisClient.PublishToRedisChannel(msg.Message.RoomUUID, bytes)
	return requestsMessage, nil
}
