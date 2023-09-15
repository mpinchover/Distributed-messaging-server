package records

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserUUID      string `validate:"required"`
	RoomUUID      string `validate:"required"`
	RoomID        int
	MessageText   string `validate:"required"`
	UUID          string
	MessageStatus string
	CreatedAtNano float64 `json:"createdAtNano"`

	SeenBy []*SeenBy
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
	UUID          string
	CreatedAtNano float64    `json:"createdAtNano"`
	Members       []*Member  `gorm:"foreignKey:RoomID;" validate:"required"`
	Messages      []*Message `gorm:"foreignKey:RoomID;" validate:"required"`
}

type Member struct {
	gorm.Model
	RoomUUID string
	RoomID   int
	UserUUID string `validate:"required"`
}

// /* AUTH   */
// // for ext service, not chat user
// type AuthProfile struct {
// 	gorm.Model
// 	UUID           string
// 	Email          string
// 	HashedPassword string
// 	Mobile         string
// }

// /* MATCHING   */

// // after user has answered
// type TrackedQuestion struct {
// 	gorm.Model
// 	UUID         string
// 	QuestionText string
// 	Category     string
// 	UserUUID     string
// 	QuestionUUID string
// 	Liked        bool
// }

// type DiscoverProfile struct {
// 	gorm.Model
// 	Gender           string
// 	GenderPreference string
// 	Age              int64
// 	MinAgePref       int64
// 	MaxAgePref       int64
// 	UserUUID         string
// 	CurrentLat       float64
// 	CurrentLng       float64
// }

// type TrackedLike struct {
// 	gorm.Model
// 	UUID       string
// 	UserUUID   string
// 	TargetUUID string
// 	Liked      bool
// }
