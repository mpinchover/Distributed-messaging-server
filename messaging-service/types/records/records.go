package records

import "gorm.io/gorm"

type ChatMessage struct {
	gorm.Model
	FromUUID    string `gorm:"from_uuid"`
	RoomUUID    string `gorm:"room_uuid"`
	MessageText string `gorm:"message_text"`
	UUID        string `gorm:"uuid"`
}

type ChatRoom struct {
	gorm.Model
	UUID         string             `gorm:"uuid"`
	Participants []*ChatParticipant `gorm:"-"`
	Messages     []*ChatMessage     `gorm:"-"`
}

type ChatParticipant struct {
	gorm.Model
	UUID     string `gorm:"uuid"`
	RoomUUID string `gorm:"room_uuid"`
	UserUUID string `gorm:"user_uuid"`
}
