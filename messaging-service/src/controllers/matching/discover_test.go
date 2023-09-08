// GetQuestionsForMatching

package matching

// import (
// 	"fmt"
// 	"testing"

// 	mockStorage "messaging-service/mocks/src/gateways/storage"
// 	mockRepo "messaging-service/mocks/src/repo"
// 	"messaging-service/src/types/requests"

// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type MatchingControllerSuite struct {
// 	suite.Suite

// 	matchingController *MatchingController
// 	mockStorage        *mockStorage.StorageInterface
// 	mockRepo           *mockRepo.RepoInterface
// }

// // this function executes before the test suite begins execution
// func (s *MatchingControllerSuite) SetupSuite() {

// }

// func (s *MatchingControllerSuite) SetupTest() {
// 	s.mockStorage = mockStorage.NewStorageInterface(s.T())
// 	s.mockRepo = mockRepo.NewRepoInterface(s.T())
// 	s.matchingController = &MatchingController{
// 		Storage: s.mockStorage,
// 		Repo:    s.mockRepo,
// 	}
// }

// func TestMatchingControllerSuite(t *testing.T) {
// 	suite.Run(t, new(MatchingControllerSuite))
// }

// func (s *MatchingControllerSuite) TestGetQuestionsForMatching() {

// 	headers := []string{"category", "text", "index", "uuid"}

// 	records := [][]string{
// 		headers,
// 	}
// 	for i := 1; i <= 50; i++ {
// 		category := fmt.Sprintf("category %d", i)
// 		question := fmt.Sprintf("question %d", i)
// 		idx := fmt.Sprintf("%d", i)
// 		uuid := fmt.Sprintf("uuid-%d", i)
// 		records = append(records, []string{category, question, idx, uuid})
// 	}

// 	s.mockStorage.On("GetQuestionsFromStorage", mock.Anything).Return(records, nil).Once()
// 	questions, err := s.matchingController.GetQuestionsForMatching(0, "key")
// 	s.NoError(err)
// 	s.NotNil(questions)
// 	// s.Len(questions, len(records)-1)

// 	for i, q := range questions {
// 		s.NotEmpty(q.Category)
// 		s.Equal(fmt.Sprintf("category %d", i+1), q.Category)
// 		s.NotEmpty(q.Text)
// 		s.Equal(fmt.Sprintf("question %d", i+1), q.Text)
// 		s.NotZero(q.Index)
// 		s.Equal(int64(i+1), q.Index)
// 		s.NotEmpty(q.UUID)
// 	}
// 	totalQuestions := []*requests.Question{}
// 	totalQuestions = append(totalQuestions, questions...)
// 	s.Len(questions, 20)
// 	s.Len(totalQuestions, 20)

// 	s.mockStorage.On("GetQuestionsFromStorage", mock.Anything).Return(records, nil).Once()
// 	questions, err = s.matchingController.GetQuestionsForMatching(int64(len(totalQuestions)), "key")
// 	s.NoError(err)
// 	s.NotNil(questions)
// 	// s.Len(questions, len(records)-1)

// 	for i, q := range questions {
// 		s.NotEmpty(q.Category)
// 		s.Equal(fmt.Sprintf("category %d", i+1+len(totalQuestions)), q.Category)
// 		s.NotEmpty(q.Text)
// 		s.Equal(fmt.Sprintf("question %d", i+1+len(totalQuestions)), q.Text)
// 		s.NotZero(q.Index)
// 		s.Equal(int64(i+1+len(totalQuestions)), q.Index)
// 		s.NotEmpty(q.UUID)
// 	}
// 	totalQuestions = append(totalQuestions, questions...)
// 	s.Len(questions, 20)
// 	s.Len(totalQuestions, 40)

// 	s.mockStorage.On("GetQuestionsFromStorage", mock.Anything).Return(records, nil).Once()
// 	questions, err = s.matchingController.GetQuestionsForMatching(int64(len(totalQuestions)), "key")
// 	s.NoError(err)
// 	s.NotNil(questions)
// 	// s.Len(questions, len(records)-1)

// 	for i, q := range questions {
// 		s.NotEmpty(q.Category)
// 		s.Equal(fmt.Sprintf("category %d", i+1+len(totalQuestions)), q.Category)
// 		s.NotEmpty(q.Text)
// 		s.Equal(fmt.Sprintf("question %d", i+1+len(totalQuestions)), q.Text)
// 		s.NotZero(q.Index)
// 		s.Equal(int64(i+1+len(totalQuestions)), q.Index)
// 		s.NotEmpty(q.UUID)
// 	}
// 	totalQuestions = append(totalQuestions, questions...)
// 	s.Len(questions, 10)
// 	s.Len(totalQuestions, 50)

// }
