package records

type ChatMessage struct {
	FromUUID    string `gorm:"from_uuid"`
	RoomUUID    string `gorm:"room_uuid"`
	MessageText string `gorm:"message_text"`
	UUID        string `gorm:"uuid"`
}

type ChatRoom struct {
	UUID         string   `gorm:"uuid"`
}

type ChatParticipant struct {
	UUID     string `gorm:"uuid"`
	RoomUUID string `gorm:"room_uuid"`
	UserUUID string `gorm:"user_uuid"`
}
