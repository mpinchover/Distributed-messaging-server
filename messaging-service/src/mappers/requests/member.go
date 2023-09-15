package mappers

import (
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
)

func ToRequestMembers(in []*records.Member) []*requests.Member {
	members := make([]*requests.Member, len(in))
	for i, val := range in {
		members[i] = ToRequestMember(val)
	}
	return members
}

func ToRequestMember(member *records.Member) *requests.Member {
	return &requests.Member{
		UserUUID: member.UserUUID,
		RoomUUID: member.RoomUUID,
	}
}

func ToRecordMembers(in []*requests.Member) []*records.Member {
	members := make([]*records.Member, len(in))
	for i, val := range in {
		members[i] = ToRecordMember(val)
	}
	return members
}

func ToRecordMember(member *requests.Member) *records.Member {
	return &records.Member{
		UserUUID: member.UserUUID,
		RoomUUID: member.RoomUUID,
	}
}
