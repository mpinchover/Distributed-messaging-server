package repo

import (
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
	"time"

	"gorm.io/gorm"
)

func (r *Repo) CreateMatchingPreferences(mp *records.MatchingPreferences) error {
	return r.DB.Create(mp).Error
}

func (r *Repo) GetMatchingPreferencesByUserUUID(userUUID string) (*records.MatchingPreferences, error) {
	res := &records.MatchingPreferences{}
	err := r.DB.Where("user_uuid = ?", userUUID).First(res).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return res, err
}

// auth profile uuid
// TODO – get number of liked vs not liked
func (r *Repo) GetTrackedQuestionsByUserUUID(userUUID string) ([]*records.TrackedQuestion, error) {
	results := []*records.TrackedQuestion{}

	err := r.DB.Where("user_uuid  = ?", userUUID).Find(&results).Error

	return results, err
}

func (r *Repo) CreateTrackedQuestion(trackedQuestion *records.TrackedQuestion) error {
	return r.DB.Create(trackedQuestion).Error
}

func (r *Repo) UpdateTrackedQuestion(trackedQuestions *records.TrackedQuestion) error {
	err := r.DB.Where("user_uuid = ?", trackedQuestions.UserUUID).
		Where("question_uuid = ?", trackedQuestions.QuestionUUID).
		Update("messages", trackedQuestions).Error
	return err
}

func (r *Repo) GetCandidatesByMatchingPreferences(matchingPrefs *records.MatchingPreferences, filters *requests.MatchingFilter) ([]string, error) {
	uuids := []string{}
	err := r.DB.Model(&records.MatchingPreferences{}).
		Where("zipcode in (?)", filters.Zipcodes).
		Where("gender = ?", matchingPrefs.GenderPreference).
		Where("gender_preference = ?", matchingPrefs.Gender).
		Where("age <= ?", matchingPrefs.MaxAgePref).
		Where("age >= ?", matchingPrefs.MinAgePref).
		Where("min_age_pref <= ", matchingPrefs.Age).
		Where("max_age_pref >= ", matchingPrefs.Age).
		Where("user_uuid not in (?)", filters.ExcludeUUIDs).
		Select("user_uuid").Find(&uuids).Error
	return uuids, err
}

func (r *Repo) GetBlockedCandidatesByUser(userUUID string) ([]string, error) {
	blockedCandidates := []string{}
	return blockedCandidates, nil
}

func (r *Repo) GetRecentlyMatchedUUIDs(uuid string, t time.Time) ([]string, error) {
	uuids := []string{}
	query := `
	SELECT DISTINCT m.user_uuid
	FROM members m
	JOIN rooms r ON m.room_uuid = r.uuid
	WHERE m.room_uuid IN (
		SELECT room_uuid
		FROM members
		WHERE user_uuid = ?
	)
	AND r.created_at < ?
	`
	err := r.DB.Raw(query, uuid, t).Scan(&uuids).Error
	return uuids, err
}

// need to do a join here to ensure that the correctly filtered users are being queried
// func (r *Repo) GetCandidateProfiles(userMatchingPrefs *records.MatchingPreferences, queryFilters *matchingTypes.MatchingFilter) ([]string, error) {

// 	userLikedQuestions, err := r.GetLikedQuestionUUIDsByUserUUID(userMatchingPrefs.UserUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// if the user hasn't liked enough questions, just abort here.
// 	if len(userLikedQuestions) < 50 {
// 		return nil, nil
// 	}

// 	// get everyone this user has blocked
// 	blockedUUIDs, err := r.GetBlockedCandidatesByUser(userMatchingPrefs.UserUUID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// get everything that appeared in the last two days
// 	recentlyMatchedUUIDs, err := r.GetRecentlyMatchedUUIDs(userMatchingPrefs.UserUUID, time.Now().Add(48*time.Hour*-1))
// 	if err != nil {
// 		return nil, err
// 	}

// 	queryFilters.ExcludeUUIDs = append(queryFilters.ExcludeUUIDs, recentlyMatchedUUIDs...)
// 	queryFilters.ExcludeUUIDs = append(queryFilters.ExcludeUUIDs, blockedUUIDs...)

// 	// get all the candidates that match the dating preferences and are excluded
// 	matchingPrefsCandidateUUIDs, err := r.GetCandidatesByMatchingPreferences(userMatchingPrefs, queryFilters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// no matches, abort
// 	if len(matchingPrefsCandidateUUIDs) == 0 {
// 		return nil, nil
// 	}

// 	// now get all the questions that the user has liked that have been liked by anyone that is also a matching acnddiate
// 	candidateLikedQuestions, err := r.GetQuestionsLikedByMatchedCandidateUUIDs(userLikedQuestions, matchingPrefsCandidateUUIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// no overlapping liked questions, abort
// 	if len(candidateLikedQuestions) == 0 {
// 		return nil, nil
// 	}

// 	// map question uuid -> list of candidate uuids who have liked this q
// 	likedQuestions := map[string][]string{}
// 	// keep track of how many questions a candidate likes the overlaps w what the user likes
// 	for _, q := range candidateLikedQuestions {
// 		likedQuestions[q.UserUUID] = append(likedQuestions[q.UserUUID], q.QuestionUUID)
// 	}

// 	sort.Slice(matchingPrefsCandidateUUIDs, func(i int, j int) bool {
// 		iUUID := matchingPrefsCandidateUUIDs[i]
// 		jUUID := matchingPrefsCandidateUUIDs[j]

// 		return len(likedQuestions[iUUID]) > len(likedQuestions[jUUID])
// 	})
// 	return matchingPrefsCandidateUUIDs, nil
// 	// TODO - keep track of the users' "feed" of people who they should be shown and store it in redis
// 	// when you are seeing what questions to choose for them, draw from the redis feed as well
// }

// return the question uuids that the user has liked
func (r *Repo) GetLikedQuestionUUIDsByUserUUID(userUUID string) ([]string, error) {
	results := []string{}
	err := r.DB.Where("user_uuid = ?", userUUID).Select("question_uuid").Find(&results).Error
	return results, err
}

func (r *Repo) GetLikedTrackedQuestionByUserUUIDAndCandidates(userUUID string, candidateUUIDs []string) ([]*records.TrackedQuestion, error) {
	results := []*records.TrackedQuestion{}
	query := `
		select * from tracked_questions tq
		where tq.question_uuid in (
			select question_uuid from tracked_questions tq 
			where tq.user_uuid = ?
			and tq.liked = true
			and tq.deleted_at is null
		)
		where tq.liked = true and 
		tq.user_uuid in (?) and tq.deleted_at is null
	`
	err := r.DB.Raw(query, userUUID, candidateUUIDs).Error
	return results, err
}

func (r *Repo) GetQuestionsLikedByMatchedCandidateUUIDs(questionUUIDs []string, candidateUUIDs []string) ([]*records.TrackedQuestion, error) {
	results := []*records.TrackedQuestion{}
	err := r.DB.Where("question_uuid in (?)", questionUUIDs).Where("user_uuid in (?)", candidateUUIDs).Error
	return results, err
}
