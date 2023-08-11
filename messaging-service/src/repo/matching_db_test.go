package repo

import (
	"fmt"
	"messaging-service/src/types/records"
	"messaging-service/src/types/requests"
	"messaging-service/src/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RepoSuite struct {
	suite.Suite

	repo   *Repo
	mainDB *gorm.DB
}

// this function executes before the test suite begins execution
func (r *RepoSuite) SetupSuite() {

	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/messaging?charset=utf8mb4&parseTime=True&loc=Local", "localhost", "3308")
	// make connection here
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		// FullSaveAssociations: true,
	})
	r.NoError(err)
	r.repo = &Repo{
		// DB: db,
	}
	r.mainDB = db

}

func (r *RepoSuite) SetupTest() {
	r.repo.DB = r.mainDB.Begin()

}

func (r *RepoSuite) TeardownSuite() {
	// r.repo.DB.Rollback()
}

// before each test exec the query
// https://stackoverflow.com/questions/38998267/how-to-execute-a-sql-file
// try to connect to db just to connect to gorm

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (r *RepoSuite) TestCreateTrackedQuestion() {
	// r.repo.DB = r.mainDB.Begin()
	defer r.repo.DB.Rollback()
	uuid := "tq-uuid-z"
	userUUID := "some-user-uuid-z"
	trackedQuestion := &records.TrackedQuestion{
		UUID:         uuid,
		QuestionText: "question-text-z",
		Category:     "some-category",
		UserUUID:     userUUID,
		Liked:        true,
	}
	err := r.repo.CreateTrackedQuestion(trackedQuestion)
	r.NoError(err)

	tqs, err := r.repo.GetTrackedQuestionsByUserUUID(userUUID)
	r.NoError(err)
	r.NotNil(tqs)
	r.Len(tqs, 1)

	r.Equal("question-text-z", trackedQuestion.QuestionText)
	r.Equal("some-category", trackedQuestion.Category)
	r.Equal(true, trackedQuestion.Liked)

	tqs, err = r.repo.GetTrackedQuestionsByUserUUID(userUUID + "z")
	r.NoError(err)

	r.NotNil(tqs)
	r.Len(tqs, 0)

	// r.repo.DB.Rollback()
}

func (r *RepoSuite) TestCreateMatchingPreferences() {
	// r.repo.DB = r.mainDB.Begin()
	defer r.repo.DB.Rollback()

	mp := &records.DiscoverProfile{
		Gender:           "MALE",
		GenderPreference: "FEMALE",
		Age:              30,
		MinAgePref:       25,
		MaxAgePref:       35,
		UserUUID:         "user-uuid-0",
	}
	err := r.repo.CreateDiscoverProfile(mp)
	r.NoError(err)

	// can use env variable to get the path of the working directory
	// and then you can just recreate the db every time
	res, err := r.repo.GetMatchingPreferencesByUserUUID("user-uuid-0")
	r.NoError(err)

	r.Equal(mp.Gender, res.Gender)
	r.Equal(mp.GenderPreference, res.GenderPreference)
	r.Equal(mp.Age, res.Age)
	r.Equal(mp.MinAgePref, res.MinAgePref)
	r.Equal(mp.MaxAgePref, res.MaxAgePref)
	r.Equal(mp.UserUUID, res.UserUUID)

	res, err = r.repo.GetMatchingPreferencesByUserUUID("user-uuid-abdef")
	r.NoError(err)
	r.Nil(res)

	res, err = r.repo.GetMatchingPreferencesByUserUUID("user-uuid-abdeddf")
	r.NoError(err)
	r.Nil(res)
}

func (r *RepoSuite) TestGetCandidateDiscoverProfile() {
	defer r.repo.DB.Rollback()

	// 11217
	candidateOneUUID := uuid.New().String()
	candidate1 := &records.DiscoverProfile{
		Gender:           "FEMALE",
		GenderPreference: "MALE",
		UserUUID:         candidateOneUUID,
		CurrentLat:       40.687995,
		CurrentLng:       -73.9820318,
	}
	err := r.repo.CreateDiscoverProfile(candidate1)
	r.NoError(err)

	// denver
	candidateTwoUUID := uuid.New().String()
	candidate1 = &records.DiscoverProfile{
		Gender:           "FEMALE",
		GenderPreference: "MALE",
		UserUUID:         candidateTwoUUID,
		CurrentLat:       39.7642224,
		CurrentLng:       -105.0199203,
	}
	err = r.repo.CreateDiscoverProfile(candidate1)
	r.NoError(err)

	// 06117
	candidateThreeUUID := uuid.New().String()
	profileLat := 41.8054284
	profileLng := -72.7391128
	maxDistanceMeters := int64(162000)
	filters := &requests.ProfileFilter{
		UserUUID:                &candidateThreeUUID,
		ProfileGender:           utils.ToStrPtr("MALE"),
		ProfileGenderPreference: utils.ToStrPtr("FEMALE"),
		ProfileLat:              &profileLat,
		ProfileLng:              &profileLng,
		MaxDistanceMeters:       &maxDistanceMeters,
	}

	profiles, err := r.repo.GetCandidateDiscoverProfile(filters)
	r.NoError(err)
	r.Len(profiles, 1)
	r.Equal(profiles[0].UserUUID, candidateOneUUID)

}

// func (r *RepoSuite) TestGetCandidatesByMatchingPreferences() {

// 	userMP := &records.MatchingPreferences{
// 		Zipcode:          "06117",
// 		Gender:           "MALE",
// 		GenderPreference: "FEMALE",
// 		Age:              30,
// 		MinAgePref:       25,
// 		MaxAgePref:       35,
// 		UserUUID:         "user-uuid-0",
// 	}
// 	err := r.repo.CreateMatchingPreferences(userMP)
// 	r.NoError(err)

// 	candidates := []*records.MatchingPreferences{
// 		{
// 			Zipcode:          "06112",
// 			Gender:           "FEMALE",
// 			GenderPreference: "FEMALE",
// 			Age:              30,
// 			MinAgePref:       25,
// 			MaxAgePref:       35,
// 			UserUUID:         "user-uuid-0",
// 		}
// 	}
// }
