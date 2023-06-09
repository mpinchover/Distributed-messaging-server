package repo

import (
	"errors"
	"fmt"
	"messaging-service/types/records"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/messaging?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func New() (*Repo, error) {
	var db *gorm.DB
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := connect()
	if err != nil {
		panic(err)
	}

	return &Repo{
		DB: db,
	}, nil
}

func (r *Repo) SaveRoom(room *records.ChatRoom) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(room).Error
		if err != nil {
			return err
		}

		err = tx.Create(room.Participants).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// handle the deleted_at situation too
// check if connection isn't there before doing this
// need to save chat participants on tx as well
func (r *Repo) SaveChatMessage(msg *records.ChatMessage) error {
	err := r.DB.Create(msg).Error
	return err
}

// TODO - optimize this query and how we're doing it.
// client should query messages per room
func (r *Repo) GetHyrdatedRoomsByUserUUID(uuid string) ([]*records.ChatRoom, error) {
	rooms, err := r.GetRoomsByUserUUID(uuid)
	if err != nil {
		return nil, err
	}

	// m := map[string]*records.ChatRoom{}
	// roomUUIDs := make([]string, len(rooms))
	// for i, room := range rooms {
	// 	roomUUIDs[i] = room.UUID
	// 	m[room.UUID] = room
	// }

	// messages, err := r.GetMessagesByRoomUUIDs(roomUUIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, msg := range messages {
	// 	roomUUID := msg.RoomUUID
	// 	m[roomUUID].Messages = append(m[roomUUID].Messages, msg)
	// }

	return rooms, err
}

func (r *Repo) GetMessagesByRoomUUIDs(roomUUIDs []string) ([]*records.ChatMessage, error) {
	if len(roomUUIDs) == 0 {
		return nil, errors.New("cannot query with empty list of uuids")
	}
	results := []*records.ChatMessage{}

	err := r.DB.Where("room_uuid in ?", roomUUIDs).Find(&results).Error
	return results, err
}

func (r *Repo) GetRoomsByUUIDs(uuids []string) ([]*records.ChatRoom, error) {
	if len(uuids) == 0 {
		return nil, errors.New("cannot query with empty list of uuids")
	}
	results := []*records.ChatRoom{}

	err := r.DB.Where("uuid in ?", uuids).Find(&results).Error
	return results, err
}

func (r *Repo) GetRoomsByUserUUID(uuid string) ([]*records.ChatRoom, error) {
	results := []*records.ChatRoom{}

	err := r.DB.Raw("SELECT * from chat_rooms cr join chat_participants cp on cp.user_uuid = ? and cp.room_uuid = cr.uuid", uuid).Scan(&results).Error
	return results, err
}

func (r *Repo) GetMessagesByUserUUID(uuid string) ([]*records.ChatMessage, error) {
	results := []*records.ChatMessage{}
	err := r.DB.Model(records.ChatMessage{UUID: uuid}).
		Find(&results).Error
	return results, err
}
