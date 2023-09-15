package mappers

import (
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
)

func ToRequestRooms(in []*records.Room) []*requests.Room {
	rooms := make([]*requests.Room, len(in))
	for i, val := range in {
		rooms[i] = ToRequestRoom(val)
	}
	return rooms
}

func ToRequestRoom(room *records.Room) *requests.Room {
	return &requests.Room{
		UUID:          room.UUID,
		CreatedAtNano: room.CreatedAtNano,
		Members:       ToRequestMembers(room.Members),
		Messages:      ToRequestMessages(room.Messages),
	}
}

func ToRecordRooms(in []*requests.Room) []*records.Room {
	rooms := make([]*records.Room, len(in))
	for i, val := range in {
		rooms[i] = ToRecordRoom(val)
	}
	return rooms
}

func ToRecordRoom(room *requests.Room) *records.Room {
	return &records.Room{
		UUID:          room.UUID,
		CreatedAtNano: room.CreatedAtNano,
		Members:       ToRecordMembers(room.Members),
		Messages:      ToRecordMessages(room.Messages),
	}
}
