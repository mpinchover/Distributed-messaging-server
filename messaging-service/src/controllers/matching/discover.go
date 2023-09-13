package matching

// import (
// 	"errors"
// 	"math"
// 	"messaging-service/src/types/dataframe"
// 	"messaging-service/src/types/requests"
// 	"strconv"
// )

// // check the redis feed for this user as well to see which matched profiles should be shown
// func (m *MatchingController) GetQuestionsForMatching(offset int64, key string) ([]*requests.Question, error) {
// 	// get all the questions
// 	// if you return no questions and are at end of them, just restart

// 	// got from aws s3
// 	// records := [][]string{}
// 	records, err := m.Storage.GetQuestionsFromStorage(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	df, err := dataframe.CreateDataframe(records)
// 	if err != nil {
// 		return nil, err
// 	}

// 	questions := make([]*requests.Question, len(records)-1)
// 	for i, row := range df.Rows() {

// 		text, ok := row["text"]
// 		if !ok {
// 			return nil, errors.New("missing text for question")
// 		}

// 		index, ok := row["index"]
// 		if !ok {
// 			return nil, errors.New("missing index for question")
// 		}

// 		category, ok := row["category"]
// 		if !ok {
// 			return nil, errors.New("missing category for question")
// 		}

// 		uuid, ok := row["uuid"]
// 		if !ok {
// 			return nil, errors.New("missing uuid for question")
// 		}

// 		idx, err := strconv.Atoi(*index)
// 		if err != nil {
// 			return nil, err
// 		}
// 		question := &records.Question{
// 			UUID:     *uuid,
// 			Text:     *text,
// 			Category: *category,
// 			Index:    int64(idx),
// 		}
// 		questions[i] = question

// 	}
// 	if int(offset) >= len(questions) {
// 		offset = 0
// 	}

// 	end := math.Min(float64(offset+20), float64(len(questions)))
// 	questions = questions[offset:int(end)]
// 	// get all the tracked questions for these questions that are being chosen
// 	// so you can see how many were liked. TODO - add it into the zip code later?
// 	return questions, nil
// }
