package mappers

import (
	"messaging-service/types/records"
	"messaging-service/types/requests"
)

func FromRecordsMessagesToRequestsMessages(msgs []*records.Message) []*requests.Message {
	requestsMessages := make([]*requests.Message, len(msgs))
	for i, m := range msgs {
		requestsMessages[i] = FromRecordsMessageToRequestMessage(m)
	}
	return requestsMessages

}

func FromRecordsMessageToRequestMessage(msg *records.Message) *requests.Message {
	return &requests.Message{
		UUID:        msg.UUID,
		FromUUID:    msg.FromUUID,
		RoomUUID:    msg.RoomUUID,
		MessageText: msg.MessageText,

		// CreatedAt:   msg.CreatedAt.Unix(),
	}
}

func FromRequestsMessagesToRecordsMessages(msgs []*requests.Message) []*records.Message {
	recordMessages := make([]*records.Message, len(msgs))
	for i, m := range msgs {
		recordMessages[i] = FromRequestsMessageToRecorsMessage(m)
	}
	return recordMessages
}

func FromRequestsMessageToRecorsMessage(msg *requests.Message) *records.Message {
	return &records.Message{
		UUID:        msg.UUID,
		FromUUID:    msg.FromUUID,
		RoomUUID:    msg.RoomUUID,
		MessageText: msg.MessageText,
	}
}
