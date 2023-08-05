package matching

import "messaging-service/src/types/requests"

func (m *MatchingController) GetQuestionsForMatching(offset int64) {
	// get all the questions
	// if you return no questions and are at end of them, just restart

}

func (m *MatchingController) UpdateQuestionResponse(question *requests.Question) {
	// fetch q from database
	// for now, just update the `liked`
}
