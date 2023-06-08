package repo

import (
	"fmt"
	"messaging-service/types/records"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func New() (*Repo, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/messaging?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repo{
		DB: db,
	}, nil
}

// handle the deleted_at situation too
// check if connection isn't there before doing this
func (r *Repo) SaveChatMessage(msg *records.ChatMessage) error {
	err := r.DB.Create(msg).Error
	return err
}

func (r *Repo) GetRoomsByUserUUID(uuid string) ([]*records.ChatRoom, error) {
	results := []*records.ChatRoom{}

	err := r.DB.
		InnerJoins("chat_participants", r.DB.
			Where(&records.ChatParticipant{UUID: uuid})).
		Find(&results).Error
	return results, err
}

func (r *Repo) GetMessagesByUserUUID(uuid string) ([]*records.ChatMessage, error) {
	results := []*records.ChatMessage{}
	err := r.DB.Model(records.ChatMessage{UUID: uuid}).
		Find(&results).Error
	return results, err
}
