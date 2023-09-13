package matching

// import (
// 	"messaging-service/src/gateways/storage"
// 	redisClient "messaging-service/src/redis"
// 	"messaging-service/src/repo"
// 	"messaging-service/src/types/enums"
// 	"messaging-service/src/types/records"
// 	"messaging-service/src/types/requests"
// 	"sort"
// 	"time"
// )

// type MatchingController struct {
// 	Repo        repo.RepoInterface
// 	RedisClient redisClient.RedisInterface
// 	Storage     storage.StorageInterface
// 	// TODO - send match result over redis channels

// }

// func New(
// 	repo *repo.Repo,
// 	redisClient *redisClient.RedisClient,
// 	storage *storage.Storage,
// ) *MatchingController {

// 	matchingController := &MatchingController{
// 		Repo:        repo,
// 		RedisClient: redisClient,
// 		Storage:     storage,
// 	}

// 	return matchingController
// }

// func (m *MatchingController) UpdateQuestionResponse(q *requests.TrackedQuestion) error {
// 	return m.Repo.UpdateTrackedQuestion(&records.TrackedQuestion{
// 		UUID:         q.UUID,
// 		QuestionText: q.Text,
// 		QuestionUUID: q.QuestionUUID,
// 		UserUUID:     q.UserUUID,
// 		Liked:        q.Liked,
// 		Category:     q.Category,
// 	})
// 	// fetch q from database
// 	// for now, just update the `liked`
// }

// func (m *MatchingController) CreateMatchingFilters(userDiscoverProfile *requests.DiscoverProfile) (*requests.ProfileFilter, error) {
// 	filters := &records.ProfileFilter{}
// 	// get everyone this user has blocked
// 	blockedUUIDs, err := m.Repo.GetBlockedCandidatesByUser(userDiscoverProfile.UserUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// get everything that appeared in the last two days
// 	recentlyMatchedUUIDs, err := m.Repo.GetRecentlyMatchedUUIDs(userDiscoverProfile.UserUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// find everyone this user has already made a decision on within the past 3 days
// 	t := time.Now().Add(time.Hour * 24 * 3 * -1)
// 	recentlyMadeDecisionOn, err := m.Repo.GetRecentTrackedLikedTargetsByUserUUID(userDiscoverProfile.UserUUID, t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// add these uuids to exclude them from matching
// 	filters.ExcludeUUIDs = append(filters.ExcludeUUIDs, blockedUUIDs...)
// 	filters.ExcludeUUIDs = append(filters.ExcludeUUIDs, recentlyMatchedUUIDs...)
// 	filters.ExcludeUUIDs = append(filters.ExcludeUUIDs, recentlyMadeDecisionOn...)

// 	filters.ProfileMaxAgePreference = &userDiscoverProfile.MaxAgePref
// 	filters.ProfileMinAgePreference = &userDiscoverProfile.MinAgePref
// 	filters.ProfileGender = &userDiscoverProfile.Gender
// 	filters.ProfileGenderPreference = &userDiscoverProfile.GenderPreference
// 	filters.ProfileAge = &userDiscoverProfile.Age
// 	filters.ProfileLng = &userDiscoverProfile.CurrentLng
// 	filters.ProfileLat = &userDiscoverProfile.CurrentLat

// 	return filters, nil
// }

// func (m *MatchingController) GetQuestionsLikedByMatchedCandidates(candidateDiscoverProfiles []*requests.DiscoverProfile, userLikedQuestions []string) ([]*records.TrackedQuestion, error) {
// 	candidateUUIDs := make([]string, len(candidateDiscoverProfiles))
// 	for i, dp := range candidateDiscoverProfiles {
// 		candidateUUIDs[i] = dp.UserUUID
// 	}

// 	// now get all the questions that the user has liked that have been liked by anyone that is also a matching candidate
// 	return m.Repo.GetQuestionsLikedByMatchedCandidateUUIDs(userLikedQuestions, candidateUUIDs)
// }

// func (m *MatchingController) GetCandidateProfiles(userProfile *requests.DiscoverProfile) ([]*requests.DiscoverProfile, error) {
// 	filters, err := m.CreateMatchingFilters(userProfile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// get all the candidates that match the user dating preferences and are not excluded
// 	recordsCandidatesDiscoverProfiles, err := m.Repo.GetCandidateDiscoverProfile(filters)
// 	if err != nil {
// 		return nil, err
// 	}
// 	candidateDiscoverProfiles := make([]*requests.DiscoverProfile, len(recordsCandidatesDiscoverProfiles))
// 	for i, dp := range recordsCandidatesDiscoverProfiles {
// 		candidateDiscoverProfiles[i] = &records.DiscoverProfile{
// 			Gender:           dp.Gender,
// 			GenderPreference: dp.GenderPreference,
// 			Age:              dp.Age,
// 			MinAgePref:       dp.MinAgePref,
// 			MaxAgePref:       dp.MaxAgePref,
// 			UserUUID:         dp.UserUUID,
// 		}
// 	}

// 	return candidateDiscoverProfiles, err
// }

// // user has swiped on a question, check to see if there are any potential matches
// func (m *MatchingController) CheckMatchesForUser(userDiscoverProfile *requests.DiscoverProfile) (*requests.MatchesForUserResult, error) {

// 	res := &records.MatchesForUserResult{}
// 	// get the questions the user has liked first
// 	userLikedQuestions, err := m.Repo.GetLikedQuestionUUIDsByUserUUID(userDiscoverProfile.UserUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// if the user hasn't liked enough questions, just abort here.
// 	// TODO, create a no-op error or some kind of struct that explains why no matches
// 	if len(userLikedQuestions) < 50 {
// 		res.AbortCode = enums.ABORT_CODE_NEED_MORE_LIKED_QUESTIONS.String()
// 		return res, nil
// 	}

// 	// get all the candidates that match the user dating preferences and are not excluded
// 	candidateProfiles, err := m.GetCandidateProfiles(userDiscoverProfile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// no matches, abort
// 	if len(candidateProfiles) == 0 {
// 		res.AbortCode = enums.ABORT_CODE_NO_MATCHES.String()
// 		return res, nil
// 	}

// 	candidateLikedQuestions, err := m.GetQuestionsLikedByMatchedCandidates(candidateProfiles, userLikedQuestions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// no overlapping liked questions, abort
// 	if len(candidateLikedQuestions) == 0 {
// 		res.AbortCode = enums.ABORT_CODE_NO_OVERLAPPING_QUESTIONS.String()
// 		return res, nil
// 	}

// 	res.CandidatesMatchingPrefs = rankAndOrderProfiles(candidateLikedQuestions, candidateProfiles)
// 	// last step should be to get the profiles of the candidates
// 	return res, nil
// 	// TODO - keep track of the users' "feed" of people who they should be shown and store it in redis
// 	// when you are seeing what questions to choose for them, draw from the redis feed as well
// }

// func rankAndOrderProfiles(candidateLikedQuestions []*records.TrackedQuestion, candidateDiscoverProfiles []*requests.DiscoverProfile) []*requests.DiscoverProfile {
// 	// map question uuid -> list of candidate uuids who have liked this q
// 	likedQuestions := map[string][]string{}

// 	// keep track of how many questions a candidate likes the overlaps w what the user likes
// 	for _, q := range candidateLikedQuestions {
// 		likedQuestions[q.UserUUID] = append(likedQuestions[q.UserUUID], q.QuestionUUID)
// 	}

// 	// sort the candidates matching the dating preferences by the freq they and the user have liked the same q
// 	sort.Slice(candidateDiscoverProfiles, func(i int, j int) bool {
// 		iUUID := candidateDiscoverProfiles[i].UserUUID
// 		jUUID := candidateDiscoverProfiles[j].UserUUID

// 		return len(likedQuestions[iUUID]) > len(likedQuestions[jUUID])
// 	})

// 	return candidateDiscoverProfiles
// }

// // is this candidate a match or not
// // if yes, return their profile
// // TODO - take care of blocks
// // func (m *MatchingController) RankCandidates(userUUID string, candidateUUID string) (*matchingTypes.Profile, error) {
// // 	commonLikedQuestions, err := m.GetCommonLikedQuestions(userUUID, candidateUUID)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// TODO - get better at determining this. Maybe just need really focused questions?
// // 	if len(commonLikedQuestions) < 10 {
// // 		return nil, nil
// // 	}

// // 	// it's a match
// // 	// save the match as a room
// // 	roomUUID := uuid.New().String()
// // 	room := &records.Room{
// // 		UUID:          roomUUID,
// // 		CreatedAtNano: float64(time.Now().UnixNano()),
// // 		Members: []*records.Member{
// // 			{
// // 				UUID:     uuid.New().String(),
// // 				RoomUUID: roomUUID,
// // 				UserUUID: userUUID,
// // 			},
// // 			{
// // 				UUID:     uuid.New().String(),
// // 				RoomUUID: roomUUID,
// // 				UserUUID: candidateUUID,
// // 			},
// // 		},
// // 	}
// // 	err = m.Repo.SaveRoom(room)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	membersAsRequest := make([]*records.Member, len(room.Members))
// // 	for i, m := range room.Members {
// // 		membersAsRequest[i] = &records.Member{
// // 			UUID:     m.UUID,
// // 			UserUUID: m.UserUUID,
// // 		}
// // 	}

// // 	roomAsRequest := &records.Room{
// // 		UUID:          room.UUID,
// // 		CreatedAtNano: room.CreatedAtNano,
// // 		Members:       membersAsRequest,
// // 	}

// // 	openRoomEvent := &requests.OpenRoomEvent{
// // 		EventType: enums.EVENT_OPEN_ROOM.String(),
// // 		Room:      roomAsRequest,
// // 	}

// // 	// send the event out over redis
// // 	bytes, err := json.Marshal(openRoomEvent)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// TODO - handle event for create room
// // 	m.RedisClient.PublishToRedisChannel(enums.CHANNEL_SERVER_EVENTS, bytes)

// // 	// return the candidate
// // 	return &matchingTypes.Profile{
// // 		UUID: candidateUUID,
// // 	}, nil
// // }

// // // you will have to filter the users as well, so just do the query in 2 parts.
// // // first get the users UUIDs, then run the filtering query
// // // TODO - optimize these queries
// // func (m *MatchingController) GetCandidatesWhoLikeSameQuestions(userUUID string) ([]*matchingTypes.Question, error) {
// // 	// get the questions this user has liked
// // 	likedTrackedQuestionsByUser, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(userUUID)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// question UUIDs
// // 	questionUUIDs := make([]string, len(likedTrackedQuestionsByUser))
// // 	for i, q := range likedTrackedQuestionsByUser {
// // 		questionUUIDs[i] = q.UUID
// // 	}
// // 	allLikedTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByQuestionUUIDs(questionUUIDs)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// now rank them
// // 	userUUIDToLikedQuestions := map[string][]*records.TrackedQuestion{}

// // 	for _, trackedQ := range allLikedTrackedQuestions {
// // 		userUUIDToLikedQuestions[trackedQ.UserUUID] = append(userUUIDToLikedQuestions[trackedQ.UserUUID], trackedQ)
// // 	}

// // 	// get all the people who have liked these questions
// // 	candidateUUIDs := make([]string, len(allLikedTrackedQuestions)-1)
// // 	for i, tq := range allLikedTrackedQuestions {
// // 		if tq.UserUUID == userUUID {
// // 			continue
// // 		}
// // 		candidateUUIDs[i] = tq.UserUUID
// // 	}
// // }

// // // TODO - show you matched on 3 + (number of other matches) things
// // func (m *MatchingController) GetCommonLikedQuestions(userUUID string, candidateUUID string) ([]*matchingTypes.Question, error) {
// // 	userTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(userUUID)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	uTrackedQuestions := make([]*matchingTypes.TrackedQuestion, len(userTrackedQuestions))
// // 	for i, q := range userTrackedQuestions {
// // 		reqQuestion := &matchingTypes.TrackedQuestion{
// // 			Text:         q.Text,
// // 			UUID:         q.UUID,
// // 			Category:     q.Category,
// // 			Index:        q.Index,
// // 			Liked:        q.Liked,
// // 			QuestionUUID: q.QuestionUUID,
// // 			UserUUID:     q.UserUUID,
// // 		}
// // 		uTrackedQuestions[i] = reqQuestion
// // 	}

// // 	// TODO - convert to requests
// // 	candidateTrackedQuestions, err := m.Repo.GetLikedTrackedQuestionsByUserUUID(candidateUUID)
// // 	cTrackedQuestions := make([]*matchingTypes.TrackedQuestion, len(candidateTrackedQuestions))
// // 	for i, q := range candidateTrackedQuestions {
// // 		reqQuestion := &matchingTypes.TrackedQuestion{
// // 			Text:         q.Text,
// // 			UUID:         q.UUID,
// // 			Category:     q.Category,
// // 			Index:        q.Index,
// // 			Liked:        q.Liked,
// // 			QuestionUUID: q.QuestionUUID,
// // 			UserUUID:     q.UserUUID,
// // 		}
// // 		cTrackedQuestions[i] = reqQuestion
// // 	}
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	candidateQuestionsMap := map[string]*matchingTypes.TrackedQuestion{}
// // 	for _, q := range cTrackedQuestions {
// // 		candidateQuestionsMap[q.UUID] = q
// // 	}

// // 	commonLikes := []*matchingTypes.Question{}
// // 	for _, q := range userTrackedQuestions {
// // 		if _, ok := candidateQuestionsMap[q.UUID]; !ok {
// // 			continue
// // 		}

// // 		// start if with just things people like
// // 		if q.Liked && q.Liked == candidateQuestionsMap[q.UUID].Liked {
// // 			question := &matchingTypes.Question{
// // 				UUID:     q.QuestionUUID,
// // 				Index:    q.Index,
// // 				Text:     q.Text,
// // 				Category: q.Category,
// // 			}
// // 			commonLikes = append(commonLikes, question)
// // 		}
// // 	}

// // 	return commonLikes, nil
// // }
