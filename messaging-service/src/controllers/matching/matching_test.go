package matching

// import (
// 	"errors"
// 	"messaging-service/src/types/enums"
// 	"messaging-service/src/types/records"
// 	"messaging-service/src/types/requests"
// 	"messaging-service/src/utils"

// 	"github.com/stretchr/testify/mock"
// )

// func (s *MatchingControllerSuite) TestCheckMatchesForUser() {
// 	tests := []struct {
// 		test        string
// 		user        *requests.DiscoverProfile
// 		mocks       func()
// 		expectedErr string
// 		expectedRes *requests.MatchesForUserResult
// 	}{
// 		{
// 			user: &records.DiscoverProfile{
// 				UserUUID: "user-uuid",
// 			},
// 			mocks: func() {
// 				s.mockRepo.On("GetLikedQuestionUUIDsByUserUUID", mock.Anything).Return(nil, errors.New("GetLikedQuestionUUIDsByUserUUID failed")).Once()
// 			},
// 			expectedErr: "GetLikedQuestionUUIDsByUserUUID failed",
// 		},
// 		{
// 			test: "not enough liked questions",
// 			user: &records.DiscoverProfile{
// 				UserUUID: "user-uuid",
// 			},
// 			mocks: func() {
// 				likedQuestionUUIDs := make([]string, 20)
// 				for i := 0; i < len(likedQuestionUUIDs); i++ {
// 					likedQuestionUUIDs[i] = "uuid"
// 				}
// 				s.mockRepo.On("GetLikedQuestionUUIDsByUserUUID", mock.Anything).Return(likedQuestionUUIDs, nil).Once()
// 			},
// 			expectedRes: &records.MatchesForUserResult{
// 				AbortCode: enums.ABORT_CODE_NEED_MORE_LIKED_QUESTIONS.String(),
// 			},
// 		},
// 		{
// 			test: "No matches returned",
// 			user: &records.DiscoverProfile{
// 				UserUUID: "user-uuid",
// 			},
// 			mocks: func() {
// 				likedQuestionUUIDs := make([]string, 50)
// 				for i := 0; i < len(likedQuestionUUIDs); i++ {
// 					likedQuestionUUIDs[i] = "uuid"
// 				}
// 				s.mockRepo.On("GetLikedQuestionUUIDsByUserUUID", mock.Anything).Return(likedQuestionUUIDs, nil).Once()
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return([]string{"user-3", "user-4"}, nil).Once()
// 				s.mockRepo.On("GetRecentTrackedLikedTargetsByUserUUID", mock.Anything, mock.Anything).Return([]string{"user-5", "user-6"}, nil).Once()
// 				s.mockRepo.On("GetCandidateDiscoverProfile", mock.Anything).Return([]*records.DiscoverProfile{}, nil).Once()
// 			},
// 			expectedRes: &records.MatchesForUserResult{
// 				AbortCode: enums.ABORT_CODE_NO_MATCHES.String(),
// 			},
// 		},
// 		{
// 			test: "No overlapping questions returned",
// 			user: &records.DiscoverProfile{
// 				UserUUID: "user-uuid",
// 			},
// 			mocks: func() {
// 				likedQuestionUUIDs := make([]string, 50)
// 				for i := 0; i < len(likedQuestionUUIDs); i++ {
// 					likedQuestionUUIDs[i] = "uuid"
// 				}

// 				candidates := make([]*records.DiscoverProfile, 50)
// 				for i := 0; i < len(candidates); i++ {
// 					candidates[i] = &records.DiscoverProfile{}
// 				}

// 				s.mockRepo.On("GetLikedQuestionUUIDsByUserUUID", mock.Anything).Return(likedQuestionUUIDs, nil).Once()
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return([]string{"user-3", "user-4"}, nil).Once()
// 				s.mockRepo.On("GetRecentTrackedLikedTargetsByUserUUID", mock.Anything, mock.Anything).Return([]string{"user-5", "user-6"}, nil).Once()
// 				s.mockRepo.On("GetCandidateDiscoverProfile", mock.Anything).Return(candidates, nil).Once()
// 				s.mockRepo.On("GetQuestionsLikedByMatchedCandidateUUIDs", mock.Anything, mock.Anything).Return([]*records.TrackedQuestion{}, nil).Once()
// 			},
// 			expectedRes: &records.MatchesForUserResult{
// 				AbortCode: enums.ABORT_CODE_NO_OVERLAPPING_QUESTIONS.String(),
// 			},
// 		},
// 		{
// 			test: "success",
// 			user: &records.DiscoverProfile{
// 				UserUUID: "user-uuid",
// 			},
// 			mocks: func() {
// 				likedQuestionUUIDs := make([]string, 50)
// 				for i := 0; i < len(likedQuestionUUIDs); i++ {
// 					likedQuestionUUIDs[i] = "uuid"
// 				}

// 				candidates := make([]*records.DiscoverProfile, 50)
// 				for i := 0; i < len(candidates); i++ {
// 					candidates[i] = &records.DiscoverProfile{}
// 				}

// 				trackedQuestions := make([]*records.TrackedQuestion, 10)

// 				s.mockRepo.On("GetLikedQuestionUUIDsByUserUUID", mock.Anything).Return(likedQuestionUUIDs, nil).Once()
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return([]string{"user-3", "user-4"}, nil).Once()
// 				s.mockRepo.On("GetRecentTrackedLikedTargetsByUserUUID", mock.Anything, mock.Anything).Return([]string{"user-5", "user-6"}, nil).Once()
// 				s.mockRepo.On("GetCandidateDiscoverProfile", mock.Anything).Return(candidates, nil).Once()
// 				s.mockRepo.On("GetQuestionsLikedByMatchedCandidateUUIDs", mock.Anything, mock.Anything).Return([]*records.TrackedQuestion{}, nil).Once()
// 			},
// 			expectedRes: &records.MatchesForUserResult{
// 				AbortCode: enums.ABORT_CODE_NO_OVERLAPPING_QUESTIONS.String(),
// 			},
// 		},
// 	}

// 	for _, t := range tests {
// 		t.mocks()

// 		res, err := s.matchingController.CheckMatchesForUser(t.user)
// 		if t.expectedErr != "" {
// 			s.Error(err)
// 			s.Contains(err.Error(), t.expectedErr)
// 		} else {
// 			s.NotNil(res)
// 			s.Equal(t.expectedRes, res)
// 		}
// 	}
// }

// func (s *MatchingControllerSuite) TestCreateMatchingFilters() {
// 	tests := []struct {
// 		test            string
// 		user            *requests.DiscoverProfile
// 		mocks           func()
// 		expectedErr     string
// 		expectedFilters *requests.ProfileFilter
// 	}{
// 		{
// 			user: &records.DiscoverProfile{},
// 			test: "GetBlockedCandidatesByUser fails",
// 			mocks: func() {
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return(nil, errors.New("GetBlockedCandidatesByUser failed")).Once()
// 			},
// 			expectedErr: "GetBlockedCandidatesByUser failed",
// 		},
// 		{
// 			user: &records.DiscoverProfile{},
// 			test: "GetBlockedCandidatesByUser fails",
// 			mocks: func() {
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return(nil, errors.New("GetRecentlyMatchedUUIDs failed")).Once()
// 			},
// 			expectedErr: "GetRecentlyMatchedUUIDs failed",
// 		},
// 		{
// 			user: &records.DiscoverProfile{},
// 			test: "GetBlockedCandidatesByUser fails",
// 			mocks: func() {
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return([]string{"user-3", "user-4"}, nil).Once()
// 				s.mockRepo.On("GetRecentTrackedLikedTargetsByUserUUID", mock.Anything, mock.Anything).Return(nil, errors.New("GetRecentTrackedLikedTargetsByUserUUID failed")).Once()
// 			},
// 			expectedErr: "GetRecentTrackedLikedTargetsByUserUUID failed",
// 		},
// 		{
// 			user: &records.DiscoverProfile{
// 				UserUUID:         "user-uuid-1",
// 				Gender:           "MALE",
// 				GenderPreference: "FEMALE",
// 				Age:              30,
// 				MinAgePref:       25,
// 				MaxAgePref:       35,
// 				CurrentLat:       10,
// 				CurrentLng:       12,
// 			},
// 			test: "GetBlockedCandidatesByUser fails",
// 			mocks: func() {
// 				s.mockRepo.On("GetBlockedCandidatesByUser", mock.Anything).Return([]string{"user-1", "user-2"}, nil).Once()
// 				s.mockRepo.On("GetRecentlyMatchedUUIDs", mock.Anything).Return([]string{"user-3", "user-4"}, nil).Once()
// 				s.mockRepo.On("GetRecentTrackedLikedTargetsByUserUUID", mock.Anything, mock.Anything).Return([]string{"user-5", "user-6"}, nil).Once()
// 			},
// 			expectedFilters: &records.ProfileFilter{
// 				ProfileGender:           utils.ToStrPtr("MALE"),
// 				ProfileGenderPreference: utils.ToStrPtr("FEMALE"),
// 				ExcludeUUIDs:            []string{"user-1", "user-2", "user-3", "user-4", "user-5", "user-6"},
// 				ProfileAge:              utils.ToInt64Ptr(30),
// 				ProfileMinAgePreference: utils.ToInt64Ptr(25),
// 				ProfileMaxAgePreference: utils.ToInt64Ptr(35),
// 				ProfileLat:              utils.ToFloat64Ptr(10),
// 				ProfileLng:              utils.ToFloat64Ptr(12),
// 			},
// 		},
// 	}

// 	for _, t := range tests {
// 		t.mocks()

// 		filters, err := s.matchingController.CreateMatchingFilters(t.user)
// 		if t.expectedErr != "" {
// 			s.Error(err)
// 			s.Contains(err.Error(), t.expectedErr)
// 		} else {
// 			s.NoError(err)
// 			s.NotNil(filters)
// 			s.Equal(t.expectedFilters, filters)
// 		}
// 	}
// }
