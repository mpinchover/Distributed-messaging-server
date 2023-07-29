package repo

import (
	"fmt"
	"messaging-service/src/types/records"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	PAGINATION_MESSAGES = 20
	PAGINATION_ROOMS    = 10
)

type RepoInteface interface {
	GetAuthProfileByEmail(email string) (*records.AuthProfile, error)
	SaveAuthProfile(authProfile *records.AuthProfile) error
	UpdatePassword(email string, hashedPassword string) error
	LeaveRoom(userUUID string, roomUUID string) error
	UpdateMessage(message *records.Message) error
	GetMembersByRoomUUID(roomUUID string) ([]*records.Member, error)
	GetMessageByUUID(uuid string) (*records.Message, error)
	SaveSeenBy(seenBy *records.SeenBy) error
	GetRoomByRoomUUID(roomUUID string) (*records.Room, error)
	SaveMessage(msg *records.Message) error
	GetMessagesByRoomUUID(roomUUID string, offset int) ([]*records.Message, error)
	GetMessagesByRoomUUIDs(roomUUIDs string, offset int) ([]*records.Message, error)
	GetRoomsByUserUUID(uuid string, offset int) ([]*records.Room, error)
	DeleteRoom(roomUUID string) error
	SaveRoom(room *records.Room) error
}

type Repo struct {
	DB *gorm.DB
}

func connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/messaging?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		// FullSaveAssociations: true,
	})
}

func New() *Repo {
	var db *gorm.DB
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := connect()
	if err != nil {
		panic(err)
	}

	return &Repo{
		DB: db,
	}
}
