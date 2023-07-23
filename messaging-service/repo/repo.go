package repo

import (
	"fmt"
	"messaging-service/types/records"
	"os"
	"sort"

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

func (r *Repo) SaveRoom(room *records.Room) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(room).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *Repo) LeaveRoom(userUUID string, roomUUID string) error {
	return r.DB.Where("user_uuid = ?", userUUID).
		Where("room_uuid = ?", roomUUID).
		Delete(&records.Member{}).Error
}

func (r *Repo) UpdateMessage(message *records.Message) error {
	err := r.DB.Where("uuid = ?", message.UUID).Update("messages", message).Error
	return err
}

func (r *Repo) GetMembersByRoomUUID(roomUUID string) ([]*records.Member, error) {
	result := []*records.Member{}
	err := r.DB.Where("room_uuid = ?", roomUUID).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) GetMessageByUUID(uuid string) (*records.Message, error) {
	result := &records.Message{}
	err := r.DB.Preload("SeenBy").Where("uuid = ?", uuid).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) SaveSeenBy(seenBy *records.SeenBy) error {
	return r.DB.Create(seenBy).Error
}

func (r *Repo) GetRoomByRoomUUID(roomUUID string) (*records.Room, error) {
	result := &records.Room{}
	err := r.DB.Preload("Messages").Preload("Messages.SeenBy").Preload("Members").Where("uuid = ?", roomUUID).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// handle the deleted_at situation too
// check if connection isn't there before doing this
// need to save chat participants on tx as well
func (r *Repo) SaveMessage(msg *records.Message) error {
	err := r.DB.Create(msg).Error
	return err
}

func (r *Repo) GetMessagesByRoomUUID(roomUUID string, offset int) ([]*records.Message, error) {
	results := []*records.Message{}

	query := r.DB.Preload("SeenBy").Where("room_uuid = ?", roomUUID).Order("id desc").Offset(offset).Limit(PAGINATION_MESSAGES)
	err := query.Find(&results).Error

	return results, err
}

func (r *Repo) GetMessagesByRoomUUIDs(roomUUIDs string, offset int) ([]*records.Message, error) {
	results := []*records.Message{}

	query := r.DB.Preload("SeenBy").Where("room_uuid in (?)", roomUUIDs)
	err := query.Find(&results).Error

	return results, err
}

// order the rooms by what the latest message is
// TODO - possibly switch to go routines to handle the timeouts
func (r *Repo) GetRoomsByUserUUID(uuid string, offset int) ([]*records.Room, error) {
	results := []*records.Room{}

	// query := `
	// 	SELECT r.id AS room_id, r.uuid AS room_uuid, r.created_at AS room_created_at, r.updated_at AS room_updated_at
	// 	FROM rooms r
	// 	INNER JOIN members m ON r.id = m.room_id
	// 	LEFT JOIN (
	// 		SELECT room_id, MAX(id) AS latest_message_id
	// 		FROM messages
	// 		GROUP BY room_id
	// 	) latest_msg ON r.id = latest_msg.room_id
	// 	LEFT JOIN messages msg ON latest_msg.latest_message_id = msg.id
	// 	WHERE m.user_uuid = ?
	// 	ORDER BY latest_msg.latest_message_id DESC
	// 	LIMIT ?,?;
	// `
	query := `
		SELECT r.id AS room_id
		FROM rooms r
		INNER JOIN members m ON r.id = m.room_id
		LEFT JOIN (
			SELECT room_id, MAX(id) AS latest_message_id
			FROM messages
			WHERE deleted_at is null
			GROUP BY room_id
		) latest_msg ON r.id = latest_msg.room_id
		LEFT JOIN messages msg ON latest_msg.latest_message_id = msg.id
		WHERE m.user_uuid = ? AND m.deleted_at is null and r.deleted_at is null
		ORDER BY latest_msg.latest_message_id DESC
		LIMIT ?,?;
	`

	// get the last message content as well
	roomIDs := []int{}
	err := r.DB.Raw(query, uuid, offset, PAGINATION_ROOMS).Scan(&roomIDs).Error
	if err != nil {
		return nil, err
	}

	// fmt.Println("ROOM IDS")
	// fmt.Println(roomIDs)
	// panic("STOP")

	// preload messages
	// https://stackoverflow.com/questions/57782293/how-to-limit-results-of-preload-of-gorm
	err = r.DB.
		// Preload("Messages", func(tx *gorm.DB) *gorm.DB {
		// 	return tx.Order("id desc")
		// }).
		Preload("Messages", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("id desc").Find(&records.Message{})
			// return tx.Raw("select * from messages order by id desc")
		}).
		Preload("Messages.SeenBy").
		Preload("Members").
		Where("id in (?)", roomIDs).Find(&results).Error

	sort.Slice(results, func(i, j int) bool {
		if len(results[i].Messages) == 0 {
			return false
		}
		if len(results[j].Messages) == 0 {
			return true
		}
		return results[i].Messages[0].ID > results[j].Messages[0].ID
	})

	// now that you know which rooms to get, run a query that gets the hydrated rooms
	return results, err
}

func (r *Repo) deleteMembersByRoomUUIDInTx(tx *gorm.DB, roomUUID string) error {
	return tx.
		Where("room_uuid = ?", roomUUID).
		Delete(&records.Member{}).
		Error
}

func (r *Repo) deleteSeenByByMessageUUIDInTx(tx *gorm.DB, messageUUIDs []string) error {
	return tx.
		Where("message_uuid in (?)", messageUUIDs).
		Delete(&records.SeenBy{}).
		Error
}

func (r *Repo) deleteMessagesByRoomUUIDInTx(tx *gorm.DB, roomUUID string) error {
	messageUUIDs := []string{}
	if err := tx.
		Where("room_uuid = ?", roomUUID).
		Model(&records.Message{}).
		Pluck("uuid", &messageUUIDs).
		Error; err != nil {
		return err
	}

	// delete all the seen_by
	if err := r.deleteSeenByByMessageUUIDInTx(tx, messageUUIDs); err != nil {
		return err
	}

	return tx.
		Where("room_uuid = (?) ", roomUUID).
		Delete(&records.Message{}).
		Error
}

func (r *Repo) deleteRoomByUUIDInTx(tx *gorm.DB, roomUUID string) error {
	return tx.
		Where("uuid = ?", roomUUID).
		Delete(&records.Room{}).Error
}

func (r *Repo) DeleteRoom(roomUUID string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.deleteMembersByRoomUUIDInTx(tx, roomUUID); err != nil {
			return err
		}

		if err := r.deleteMessagesByRoomUUIDInTx(tx, roomUUID); err != nil {
			return err
		}

		return r.deleteRoomByUUIDInTx(tx, roomUUID)
	})
}
