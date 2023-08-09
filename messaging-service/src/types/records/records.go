package records

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromUUID      string
	RoomUUID      string
	RoomID        int
	MessageText   string
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
	Members       []*Member  `gorm:"foreignKey:RoomID;"`
	Messages      []*Message `gorm:"foreignKey:RoomID;"`
}

type Member struct {
	gorm.Model
	UUID     string
	RoomUUID string
	RoomID   int
	UserUUID string
}

/* AUTH   */
// for ext service, not chat user
type AuthProfile struct {
	gorm.Model
	UUID           string
	Email          string
	HashedPassword string
}

/* MATCHING   */

// after user has answered
type TrackedQuestion struct {
	gorm.Model
	UUID         string
	QuestionText string
	Category     string
	UserUUID     string
	QuestionUUID string
	Liked        bool
}

type MatchingPreferences struct {
	gorm.Model
	Zipcode          string
	Gender           string
	GenderPreference string
	Age              int64
	MinAgePref       int64
	MaxAgePref       int64
	UserUUID         string
}
