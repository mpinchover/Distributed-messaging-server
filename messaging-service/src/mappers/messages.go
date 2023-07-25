package mappers

import (
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
)

func FromRecordsMessagesToRequestsMessages(msgs []*records.Message) []*requests.Message {
	requestsMessages := make([]*requests.Message, len(msgs))
	for i, m := range msgs {
		requestsMessages[i] = FromRecordsMessageToRequestMessage(m)
	}
	return requestsMessages

}

func FromRecordsMessageToRequestMessage(msg *records.Message) *requests.Message {
	if msg.MessageStatus == enums.MESSAGE_STATUS_DELETED.String() {
		msg.MessageText = ""
	}
	return &requests.Message{
		UUID:        msg.UUID,
		FromUUID:    msg.FromUUID,
		RoomUUID:    msg.RoomUUID,
		MessageText: msg.MessageText,
		CreatedAt:   msg.CreatedAt.Unix(),
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
