package repo

import "messaging-service/src/types/records"

// auth profile uuid
// TODO – get number of liked vs not liked
func (r *Repo) GetTrackedQuestionsByUserUUID(userUUID string) ([]*records.TrackedQuestion, error) {
	results := []*records.TrackedQuestion{}

	err := r.DB.Where("user_uuid  = ?", userUUID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, err
}
