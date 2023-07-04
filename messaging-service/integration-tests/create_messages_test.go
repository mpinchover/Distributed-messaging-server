package integrationtests

import (
	"log"
	"time"

	"messaging-service/types/requests"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	SocketURL = "ws://localhost:9090/ws"
)

func TestRoomAndMessagesPagination(t *testing.T) {

	t.Run("test rooms and messages pagination", func(t *testing.T) {
		log.Printf("Running %s", t.Name())
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()
		cUUID := uuid.New().String()
		dUUID := uuid.New().String()

		aResp, aWebWS := setupClientConnection(t, aUUID)
		bResp, bWebWS := setupClientConnection(t, bUUID)
		cResp, cWebWS := setupClientConnection(t, cUUID)
		dResp, dWebWS := setupClientConnection(t, dUUID)

		aConnectionUUID := aResp.ConnectionUUID
		bConnectionUUID := bResp.ConnectionUUID
		cConnectionUUID := cResp.ConnectionUUID
		dConnectionUUID := dResp.ConnectionUUID

		createRoomRequest1 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
				},
				{
					UserUUID: bUUID,
				},
			},
		}
		err := openRoom(createRoomRequest1)
		assert.NoError(t, err)

		openRoomRes1 := readOpenRoomResponse(t, aWebWS, 2)
		openRoomRes1 = readOpenRoomResponse(t, bWebWS, 2)
		roomUUID1 := openRoomRes1.Room.UUID
		_ = roomUUID1

		createRoomRequest2 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
				},
				{
					UserUUID: cUUID,
				},
			},
		}
		err = openRoom(createRoomRequest2)
		assert.NoError(t, err)
		openRoomRes2 := readOpenRoomResponse(t, cWebWS, 2)
		openRoomRes2 = readOpenRoomResponse(t, aWebWS, 2)

		roomUUID2 := openRoomRes2.Room.UUID

		// send messages between A and B
		sendMessages(t, aUUID, aConnectionUUID, roomUUID1, aWebWS)
		sendMessages(t, bUUID, bConnectionUUID, roomUUID1, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		time.Sleep(1 * time.Second)
		queryMessages(t, bUUID, roomUUID1, 1)
		queryMessages(t, aUUID, roomUUID1, 2)

		// send messages between A and C
		sendMessages(t, aUUID, aConnectionUUID, roomUUID2, aWebWS)
		sendMessages(t, cUUID, cConnectionUUID, roomUUID2, cWebWS)

		recvMessages(t, aWebWS)
		recvMessages(t, cWebWS)

		time.Sleep(1 * time.Second)
		queryMessages(t, aUUID, roomUUID2, 2)
		queryMessages(t, cUUID, roomUUID2, 1)

		// create room between A and D
		createRoomReq3 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
				},
				{
					UserUUID: dUUID,
				},
			},
		}
		err = openRoom(createRoomReq3)
		assert.NoError(t, err)

		openRoomRes3 := readOpenRoomResponse(t, dWebWS, 2)
		openRoomRes3 = readOpenRoomResponse(t, aWebWS, 2)
		roomUUID3 := openRoomRes3.Room.UUID

		// send messages between A and D
		sendMessages(t, aUUID, aConnectionUUID, roomUUID3, aWebWS)
		sendMessages(t, dUUID, dConnectionUUID, roomUUID3, dWebWS)

		recvMessages(t, aWebWS)
		recvMessages(t, dWebWS)

		time.Sleep(1 * time.Second)
		queryMessages(t, aUUID, roomUUID3, 3)
		queryMessages(t, dUUID, roomUUID3, 1)

		// create room between B and C
		openRoomReq4 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: bUUID,
				},
				{
					UserUUID: cUUID,
				},
			},
		}

		err = openRoom(openRoomReq4)
		assert.NoError(t, err)
		openRoomRes4 := readOpenRoomResponse(t, bWebWS, 2)
		openRoomRes4 = readOpenRoomResponse(t, cWebWS, 2)
		roomUUID4 := openRoomRes4.Room.UUID

		// send messages between B and C
		sendMessages(t, bUUID, bConnectionUUID, roomUUID4, bWebWS)
		sendMessages(t, cUUID, cConnectionUUID, roomUUID4, cWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, cWebWS)

		time.Sleep(1 * time.Second)
		queryMessages(t, bUUID, roomUUID4, 2)
		queryMessages(t, cUUID, roomUUID4, 2)

		// create room between B and D
		openRoomRequest5 := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: bUUID,
				},
				{
					UserUUID: dUUID,
				},
			},
		}
		err = openRoom(openRoomRequest5)
		assert.NoError(t, err)

		openRoomRes5 := readOpenRoomResponse(t, dWebWS, 2)
		openRoomRes5 = readOpenRoomResponse(t, bWebWS, 2)

		// the mobiel device will get the open room msg as well
		roomUUID5 := openRoomRes5.Room.UUID

		// send messages between B and D
		sendMessages(t, bUUID, bConnectionUUID, roomUUID5, bWebWS)
		sendMessages(t, dUUID, dConnectionUUID, roomUUID5, dWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, dWebWS)

		time.Sleep(100 * time.Millisecond)
		queryMessages(t, bUUID, roomUUID5, 3)
		queryMessages(t, dUUID, roomUUID5, 2)

	})
}

// Need to get the room id first and pass it to the text message id
func TestAllConnectionsRcvMessages(t *testing.T) {
	t.Run("test all connections get msgs", func(t *testing.T) {
		log.Printf("Running test %s", t.Name())
		aUUID := uuid.New().String()
		bUUID := uuid.New().String()

		aWebResp, aWebWS := setupClientConnection(t, aUUID)
		bWebResp, bWebWS := setupClientConnection(t, bUUID)
		_, bMobileWS := setupClientConnection(t, bUUID)

		aWebConnUUID := aWebResp.ConnectionUUID
		bWebConnUUID := bWebResp.ConnectionUUID

		openRoomEvent := &requests.CreateRoomRequest{
			Members: []*requests.Member{
				{
					UserUUID: aUUID,
				},
				{
					UserUUID: bUUID,
				},
			},
		}
		err := openRoom(openRoomEvent)
		assert.NoError(t, err)

		openRoomRes := readOpenRoomResponse(t, aWebWS, 2)
		readOpenRoomResponse(t, bWebWS, 2)
		readOpenRoomResponse(t, bMobileWS, 2)
		roomUUID := openRoomRes.Room.UUID

		sendMessages(t, bUUID, aWebConnUUID, roomUUID, aWebWS)
		sendMessages(t, bUUID, bWebConnUUID, roomUUID, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		// need to recv double the msgs
		recvMessages(t, bMobileWS)
		recvMessages(t, bMobileWS)
		queryMessages(t, bUUID, roomUUID, 1)
		queryMessages(t, aUUID, roomUUID, 1)

		// add new connection
		_, aMobileWS := setupClientConnection(t, aUUID)

		sendMessages(t, bUUID, aWebConnUUID, roomUUID, aWebWS)
		sendMessages(t, bUUID, bWebConnUUID, roomUUID, bWebWS)

		recvMessages(t, bWebWS)
		recvMessages(t, aWebWS)

		// need to recv double the msgs
		recvMessages(t, bMobileWS)
		recvMessages(t, bMobileWS)

		// need to recv double the msgs
		recvMessages(t, aMobileWS)
		recvMessages(t, aMobileWS)

	})
}
