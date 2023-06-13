package repo

import (
	"fmt"
	"messaging-service/types/records"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	PAGINATION_MESSAGES = 20
	PAGINATION_ROOMS    = 10
)

type Repo struct {
	DB *gorm.DB
}

func connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/messaging?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
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

func (r *Repo) GetMessagesByRoomUUID(roomUUID string, offset int) ([]*records.ChatMessage, error) {
	results := []*records.ChatMessage{}

	query := r.DB.Where("room_uuid = ?", roomUUID).Order("id desc").Offset(offset).Limit(PAGINATION_MESSAGES)
	err := query.Find(&results).Error
	return results, err
}

func (r *Repo) GetRoomsByUserUUID(uuid string, offset int) ([]*records.ChatRoom, error) {
	results := []*records.ChatRoom{}

	// log.Println("GETTINR ROOMS")
	err := r.DB.Raw("SELECT * from chat_rooms cr join chat_participants cp on cp.user_uuid = ? and cp.room_uuid = cr.uuid limit ? offset ? ", uuid, PAGINATION_ROOMS, offset).Scan(&results).Error
	return results, err
}
