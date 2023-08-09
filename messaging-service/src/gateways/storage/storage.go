package storage

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// https://developers.google.com/sheets/api/quickstart/go

type Creds struct {
}

func (c Creds) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}, nil
}

type StorageInterface interface {
	GetQuestionsFromStorage(key string) ([][]string, error)
}

type Storage struct {
	client *s3.Client
}

func NewStorage() *Storage {
	client := s3.New(s3.Options{
		Region:      "us-east-2",
		Credentials: Creds{},
	})

	return &Storage{
		client: client,
	}
}

// // create csv from s3 file
// func createCSV(records [][]string) map[string][]*string {
// 	// if len = 1, just headers and no data
// 	if len(records) <= 1 {
// 		return nil
// 	}
// 	questions := map[string][]*string{}

// 	maxLen := 0
// 	for i := range records {
// 		maxLen = int(math.Max(float64(len(records[i])), float64(maxLen)))
// 	}

// 	indexToCol := make(map[int]string, len(records[0]))

// 	for i, col := range records[0] {
// 		questions[col] = make([]*string, maxLen)
// 		indexToCol[i] = col
// 	}

// 	for i, row := range records[1:] {
// 		for j, data := range row {
// 			col := indexToCol[j]
// 			questions[col][i] = &data
// 		}
// 	}
// 	return questions
// }

func (s *Storage) GetQuestionsFromStorage(key string) ([][]string, error) {
	rawObject, err := s.client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET")),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(rawObject.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	return records, nil

	// newCSV, err := typeCsv.CreateCSV(records)
	// // csvDF := createCSV(records)
	// return newCSV, nil
}
