package dataframe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDataframe(t *testing.T) {
	headers := []string{"header_1", "header_2", "header_3"}
	records := [][]string{
		headers,
		{"row_0_col_0", "row_0_col_1", "row_0_col_2"},
		{"row_1_col_0", "row_1_col_1", "row_1_col_2"},
		{"row_2_col_0", "row_2_col_1", "row_2_col_2"},
		{"row_3_col_0", "row_3_col_1", "row_3_col_2"},
	}
	df, err := CreateDataframe(records)
	assert.NoError(t, err)
	assert.Len(t, df.rows, 4)

	for _, row := range df.Rows() {
		for _, header := range headers {
			assert.NotNil(t, row[header], fmt.Sprintf("could not get %s", header))
		}
	}

	records = [][]string{
		headers,
		{"row_0_col_0", "", "row_0_col_2"},
		{"row_1_col_0", "", "row_1_col_2"},
		{"row_2_col_0", "", "row_2_col_2"},
		{"row_3_col_0", "", "row_3_col_2"},
	}
	df, err = CreateDataframe(records)
	assert.NoError(t, err)
	assert.Len(t, df.rows, 4)

	for _, row := range df.Rows() {
		for _, header := range headers {
			if header == "header_2" {
				assert.Nil(t, row[header])
			} else {
				assert.NotNil(t, row[header], fmt.Sprintf("could not get %s", header))
			}

		}
	}

	records = [][]string{
		{"category", "text", "index", "uuid"},
		{"category 1", "question 1", "1", "uuid-1"},
		{"category 2", "question 2", "2", "uuid-2"},
		{"category 3", "question 3", "3", "uuid-3"},
		{"category 4", "question 4", "4", "uuid-4"},
	}
	df, err = CreateDataframe(records)
	assert.NoError(t, err)
	assert.Len(t, df.rows, 4)

	for _, row := range df.Rows() {

		text, ok := row["text"]
		assert.True(t, ok)
		assert.NotNil(t, text)

		index, ok := row["index"]
		assert.True(t, ok)
		assert.NotNil(t, index)

		category, ok := row["category"]
		assert.True(t, ok)
		assert.NotNil(t, category)

		uuid, ok := row["uuid"]
		assert.True(t, ok)
		assert.NotNil(t, uuid)
	}

}
