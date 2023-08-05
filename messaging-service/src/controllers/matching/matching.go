package matching

import (
	"encoding/json"
	redisClient "messaging-service/src/redis"
	"messaging-service/src/repo"
	"messaging-service/src/types/enums"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
	"time"

	"github.com/google/uuid"
)

type MatchingController struct {
	Repo        repo.RepoInterface
	RedisClient redisClient.RedisInterface
	// TODO - send match result over redis channels

}

func New(
	repo *repo.Repo,
	redisClient *redisClient.RedisClient,
) *MatchingController {

	matchingController := &MatchingController{
		Repo:        repo,
		RedisClient: redisClient,
	}

	return matchingController
}

// is this candidate a match or not
// if yes, return their profile
// TODO - take care of blocks
func (m *MatchingController) RankCandidates(userUUID string, candidateUUID string) (*requests.Profile, error) {
	commonLikedQuestions, err := m.GetCommonLikedQuestions(userUUID, candidateUUID)
	if err != nil {
		return nil, err
	}

	// TODO - get better at determining this. Maybe just need really focused questions?
	if len(commonLikedQuestions) < 10 {
		return nil, nil
	}

	// it's a match
	// save the match as a room
	roomUUID := uuid.New().String()
	room := &records.Room{
		UUID:          roomUUID,
		CreatedAtNano: float64(time.Now().UnixNano()),
		Members: []*records.Member{
			{
				UUID:     uuid.New().String(),
				RoomUUID: roomUUID,
				UserUUID: userUUID,
			},
			{
				UUID:     uuid.New().String(),
				RoomUUID: roomUUID,
				UserUUID: candidateUUID,
			},
		},
	}
	err = m.Repo.SaveRoom(room)
	if err != nil {
		return nil, err
	}

	membersAsRequest := make([]*requests.Member, len(room.Members))
	for i, m := range room.Members {
		membersAsRequest[i] = &requests.Member{
			UUID:     m.UUID,
			UserUUID: m.UserUUID,
		}
	}

	roomAsRequest := &requests.Room{
		UUID:          room.UUID,
		CreatedAtNano: room.CreatedAtNano,
		Members:       membersAsRequest,
	}

	openRoomEvent := &requests.OpenRoomEvent{
		EventType: enums.EVENT_OPEN_ROOM.String(),
		Room:      roomAsRequest,
	}

	// send the event out over redis
	bytes, err := json.Marshal(openRoomEvent)
	if err != nil {
		return nil, err
	}

	// TODO - handle event for create room
	m.RedisClient.PublishToRedisChannel(enums.CHANNEL_SERVER_EVENTS, bytes)

	// return the candidate
	return &requests.Profile{
		UUID: candidateUUID,
	}, nil
}

// TODO - show you matched on 3 + (number of other matches) things
func (m *MatchingController) GetCommonLikedQuestions(userUUID string, candidateUUID string) ([]*requests.Question, error) {
	userTrackedQuestions, err := m.Repo.GetTrackedQuestionsByUserUUID(userUUID)
	if err != nil {
		return nil, err
	}

	uTrackedQuestions := make([]*requests.TrackedQuestion, len(userTrackedQuestions))
	for i, q := range userTrackedQuestions {
		reqQuestion := &requests.TrackedQuestion{
			Text:         q.Text,
			UUID:         q.UUID,
			Category:     q.Category,
			Index:        q.Index,
			Liked:        q.Liked,
			QuestionUUID: q.QuestionUUID,
			UserUUID:     q.UserUUID,
		}
		uTrackedQuestions[i] = reqQuestion
	}

	// TODO - convert to requests
	candidateTrackedQuestions, err := m.Repo.GetTrackedQuestionsByUserUUID(candidateUUID)
	cTrackedQuestions := make([]*requests.TrackedQuestion, len(candidateTrackedQuestions))
	for i, q := range candidateTrackedQuestions {
		reqQuestion := &requests.TrackedQuestion{
			Text:         q.Text,
			UUID:         q.UUID,
			Category:     q.Category,
			Index:        q.Index,
			Liked:        q.Liked,
			QuestionUUID: q.QuestionUUID,
			UserUUID:     q.UserUUID,
		}
		cTrackedQuestions[i] = reqQuestion
	}
	if err != nil {
		return nil, err
	}

	candidateQuestionsMap := map[string]*requests.TrackedQuestion{}
	for _, q := range cTrackedQuestions {
		candidateQuestionsMap[q.UUID] = q
	}

	commonLikes := []*requests.Question{}
	for _, q := range userTrackedQuestions {
		if _, ok := candidateQuestionsMap[q.UUID]; !ok {
			continue
		}

		// start if with just things people like
		if q.Liked && q.Liked == candidateQuestionsMap[q.UUID].Liked {
			question := &requests.Question{
				UUID:     q.QuestionUUID,
				Index:    q.Index,
				Text:     q.Text,
				Category: q.Category,
			}
			commonLikes = append(commonLikes, question)
		}
	}

	return commonLikes, nil
}
