package matching

import (
	"messaging-service/src/gateways/storage"
	redisClient "messaging-service/src/redis"
	"messaging-service/src/repo"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
	"sort"
	"time"
)

type MatchingController struct {
	Repo        repo.RepoInterface
	RedisClient redisClient.RedisInterface
	Storage     storage.StorageInterface
	// TODO - send match result over redis channels

}

func New(
	repo *repo.Repo,
	redisClient *redisClient.RedisClient,
	storage *storage.Storage,
) *MatchingController {

	matchingController := &MatchingController{
		Repo:        repo,
		RedisClient: redisClient,
		Storage:     storage,
	}

	return matchingController
}

func (m *MatchingController) UpdateQuestionResponse(q *requests.TrackedQuestion) error {
	return m.Repo.UpdateTrackedQuestion(&records.TrackedQuestion{
		UUID:         q.UUID,
		QuestionText: q.Text,
		QuestionUUID: q.QuestionUUID,
		UserUUID:     q.UserUUID,
		Liked:        q.Liked,
		Category:     q.Category,
	})
	// fetch q from database
	// for now, just update the `liked`
}

// user has swiped on a question, check to see if there are any potential matches
func (m *MatchingController) CheckMatchesForUser(userMatchingPrefs *requests.MatchingPreferences, filters *requests.MatchingFilter) ([]string, error) {

	// get the questions the user has liked first
	userLikedQuestions, err := m.Repo.GetLikedQuestionUUIDsByUserUUID(userMatchingPrefs.UserUUID)
	if err != nil {
		return nil, err
	}

	// if the user hasn't liked enough questions, just abort here.
	if len(userLikedQuestions) < 50 {
		return nil, nil
	}

	// get everyone this user has blocked
	blockedUUIDs, err := m.Repo.GetBlockedCandidatesByUser(userMatchingPrefs.UserUUID)
	if err != nil {
		return nil, err
	}

	// get everything that appeared in the last two days
	recentlyMatchedUUIDs, err := m.Repo.GetRecentlyMatchedUUIDs(userMatchingPrefs.UserUUID, time.Now().Add(48*time.Hour*-1))
	if err != nil {
		return nil, err
	}

	// add these uuids to exclude them from matching
	filters.ExcludeUUIDs = append(filters.ExcludeUUIDs, recentlyMatchedUUIDs...)
	filters.ExcludeUUIDs = append(filters.ExcludeUUIDs, blockedUUIDs...)

	userMatchPrefsAsRecord := &records.MatchingPreferences{
		Zipcode: userMatchingPrefs.Zipcode, // TODO get list of zipcodes or use lat/lng
	}

	// get all the candidates that match the user dating preferences and are not excluded
	candidatesMatchingDatingPrefs, err := m.Repo.GetCandidatesByMatchingPreferences(userMatchPrefsAsRecord, filters)
	if err != nil {
		return nil, err
	}

	// no matches, abort
	if len(candidatesMatchingDatingPrefs) == 0 {
		return nil, nil
	}

	// now get all the questions that the user has liked that have been liked by anyone that is also a matching candidate
	candidateLikedQuestions, err := m.Repo.GetQuestionsLikedByMatchedCandidateUUIDs(userLikedQuestions, candidatesMatchingDatingPrefs)
	if err != nil {
		return nil, err
	}

	// no overlapping liked questions, abort
	if len(candidateLikedQuestions) == 0 {
		return nil, nil
	}

	// map question uuid -> list of candidate uuids who have liked this q
	likedQuestions := map[string][]string{}

	// keep track of how many questions a candidate likes the overlaps w what the user likes
	for _, q := range candidateLikedQuestions {
		likedQuestions[q.UserUUID] = append(likedQuestions[q.UserUUID], q.QuestionUUID)
	}

	// sort the candidates matching the dating preferences by the freq they and the user have liked the same q
	sort.Slice(candidatesMatchingDatingPrefs, func(i int, j int) bool {
		iUUID := candidatesMatchingDatingPrefs[i]
		jUUID := candidatesMatchingDatingPrefs[j]

		return len(likedQuestions[iUUID]) > len(likedQuestions[jUUID])
	})
	return candidatesMatchingDatingPrefs, nil
	// TODO - keep track of the users' "feed" of people who they should be shown and store it in redis
	// when you are seeing what questions to choose for them, draw from the redis feed as well
}

// is this candidate a match or not
// if yes, return their profile
// TODO - take care of blocks
// func (m *MatchingController) RankCandidates(userUUID string, candidateUUID string) (*matchingTypes.Profile, error) {
// 	commonLikedQuestions, err := m.GetCommonLikedQuestions(userUUID, candidateUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TODO - get better at determining this. Maybe just need really focused questions?
// 	if len(commonLikedQuestions) < 10 {
// 		return nil, nil
// 	}

// 	// it's a match
// 	// save the match as a room
// 	roomUUID := uuid.New().String()
// 	room := &records.Room{
// 		UUID:          roomUUID,
// 		CreatedAtNano: float64(time.Now().UnixNano()),
// 		Members: []*records.Member{
// 			{
// 				UUID:     uuid.New().String(),
// 				RoomUUID: roomUUID,
// 				UserUUID: userUUID,
// 			},
// 			{
// 				UUID:     uuid.New().String(),
// 				RoomUUID: roomUUID,
// 				UserUUID: candidateUUID,
// 			},
// 		},
// 	}
// 	err = m.Repo.SaveRoom(room)
// 	if err != nil {
// 		return nil, err
// 	}

// 	membersAsRequest := make([]*requests.Member, len(room.Members))
// 	for i, m := range room.Members {
// 		membersAsRequest[i] = &requests.Member{
// 			UUID:     m.UUID,
// 			UserUUID: m.UserUUID,
// 		}
// 	}

// 	roomAsRequest := &requests.Room{
// 		UUID:          room.UUID,
// 		CreatedAtNano: room.CreatedAtNano,
// 		Members:       membersAsRequest,
// 	}

// 	openRoomEvent := &requests.OpenRoomEvent{
// 		EventType: enums.EVENT_OPEN_ROOM.String(),
// 		Room:      roomAsRequest,
// 	}

// 	// send the event out over redis
// 	bytes, err := json.Marshal(openRoomEvent)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TODO - handle event for create room
// 	m.RedisClient.PublishToRedisChannel(enums.CHANNEL_SERVER_EVENTS, bytes)

// 	// return the candidate
// 	return &matchingTypes.Profile{
// 		UUID: candidateUUID,
// 	}, nil
// }

// // you will have to filter the users as well, so just do the query in 2 parts.
// // first get the users UUIDs, then run the filtering query
// // TODO - optimize these queries
// func (m *MatchingController) GetCandidatesWhoLikeSameQuestions(userUUID string) ([]*matchingTypes.Question, error) {
// 	// get the questions this user has liked
// 	likedTrackedQuestionsByUser, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(userUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// question UUIDs
// 	questionUUIDs := make([]string, len(likedTrackedQuestionsByUser))
// 	for i, q := range likedTrackedQuestionsByUser {
// 		questionUUIDs[i] = q.UUID
// 	}
// 	allLikedTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByQuestionUUIDs(questionUUIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// now rank them
// 	userUUIDToLikedQuestions := map[string][]*records.TrackedQuestion{}

// 	for _, trackedQ := range allLikedTrackedQuestions {
// 		userUUIDToLikedQuestions[trackedQ.UserUUID] = append(userUUIDToLikedQuestions[trackedQ.UserUUID], trackedQ)
// 	}

// 	// get all the people who have liked these questions
// 	candidateUUIDs := make([]string, len(allLikedTrackedQuestions)-1)
// 	for i, tq := range allLikedTrackedQuestions {
// 		if tq.UserUUID == userUUID {
// 			continue
// 		}
// 		candidateUUIDs[i] = tq.UserUUID
// 	}
// }

// // TODO - show you matched on 3 + (number of other matches) things
// func (m *MatchingController) GetCommonLikedQuestions(userUUID string, candidateUUID string) ([]*matchingTypes.Question, error) {
// 	userTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(userUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	uTrackedQuestions := make([]*matchingTypes.TrackedQuestion, len(userTrackedQuestions))
// 	for i, q := range userTrackedQuestions {
// 		reqQuestion := &matchingTypes.TrackedQuestion{
// 			Text:         q.Text,
// 			UUID:         q.UUID,
// 			Category:     q.Category,
// 			Index:        q.Index,
// 			Liked:        q.Liked,
// 			QuestionUUID: q.QuestionUUID,
// 			UserUUID:     q.UserUUID,
// 		}
// 		uTrackedQuestions[i] = reqQuestion
// 	}

// 	// TODO - convert to requests
// 	candidateTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(candidateUUID)
// 	cTrackedQuestions := make([]*matchingTypes.TrackedQuestion, len(candidateTrackedQuestions))
// 	for i, q := range candidateTrackedQuestions {
// 		reqQuestion := &matchingTypes.TrackedQuestion{
// 			Text:         q.Text,
// 			UUID:         q.UUID,
// 			Category:     q.Category,
// 			Index:        q.Index,
// 			Liked:        q.Liked,
// 			QuestionUUID: q.QuestionUUID,
// 			UserUUID:     q.UserUUID,
// 		}
// 		cTrackedQuestions[i] = reqQuestion
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	candidateQuestionsMap := map[string]*matchingTypes.TrackedQuestion{}
// 	for _, q := range cTrackedQuestions {
// 		candidateQuestionsMap[q.UUID] = q
// 	}

// 	commonLikes := []*matchingTypes.Question{}
// 	for _, q := range userTrackedQuestions {
// 		if _, ok := candidateQuestionsMap[q.UUID]; !ok {
// 			continue
// 		}

// 		// start if with just things people like
// 		if q.Liked && q.Liked == candidateQuestionsMap[q.UUID].Liked {
// 			question := &matchingTypes.Question{
// 				UUID:     q.QuestionUUID,
// 				Index:    q.Index,
// 				Text:     q.Text,
// 				Category: q.Category,
// 			}
// 			commonLikes = append(commonLikes, question)
// 		}
// 	}

// 	return commonLikes, nil
// }
