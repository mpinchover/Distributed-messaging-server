package records

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUUID    string
	RoomUUID    string
	RoomID      int
	MessageText string
	UUID        string
}

type Room struct {
	gorm.Model
	UUID     string
	Members  []*Member  `gorm:"foreignKey:RoomID"`
	Messages []*Message `gorm:"foreignKey:RoomID"`
}

type Member struct {
	gorm.Model
	UUID     string
	RoomUUID string
	RoomID   int
	UserUUID string
	UserRole string
}
