package mappers

import (
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
)

func ToRequestSeenBys(seenBys []*records.SeenBy) []*requests.SeenBy {
	sbs := make([]*requests.SeenBy, len(seenBys))
	for i, val := range seenBys {
		sbs[i] = ToRequestSeenBy(val)
	}
	return sbs
}

func ToRequestSeenBy(seenBy *records.SeenBy) *requests.SeenBy {
	return &requests.SeenBy{
		MessageUUID: seenBy.MessageUUID,
		UserUUID:    seenBy.UserUUID,
	}
}

func ToRecordSeenBys(seenBys []*requests.SeenBy) []*records.SeenBy {
	sbs := make([]*records.SeenBy, len(seenBys))
	for i, val := range seenBys {
		sbs[i] = ToRecordSeenBy(val)
	}
	return sbs
}

func ToRecordSeenBy(seenBy *requests.SeenBy) *records.SeenBy {
	return &records.SeenBy{
		MessageUUID: seenBy.MessageUUID,
		UserUUID:    seenBy.UserUUID,
	}
}
