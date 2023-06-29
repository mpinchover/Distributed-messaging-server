package records

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromUUID    string
	RoomUUID    string
	RoomID      int
	MessageText string
	UUID        string
	SeenBy      []*SeenBy
}

type SeenBy struct {
	gorm.Model
	MessageUUID string
	UserUUID    string
	MessageID   int
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by SeenBy to `seen_by`
func (SeenBy) TableName() string {
	return "seen_by"
}

type Room struct {
	gorm.Model
	UUID     string
	Members  []*Member  `gorm:"foreignKey:RoomID;"`
	Messages []*Message `gorm:"foreignKey:RoomID;"`
}

// func (r *Room) BeforeDelete(tx *gorm.DB) error {
// 	err := tx.Where("room_uuid = ? ", r.UUID).Delete(&Message{}).Error
// 	if err != nil {
// 		return err
// 	}

// 	err = tx.Where("room_uuid = ? ", r.UUID).Delete(&Member{}).Error
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

type Member struct {
	gorm.Model
	UUID     string
	RoomUUID string
	RoomID   int
	UserUUID string
	UserRole string
}
