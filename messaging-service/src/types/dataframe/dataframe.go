package dataframe

import (
	"errors"
)

type Dataframe struct {
	headers []string
	data    [][]string

	rows []map[string]*string
}

func (c *Dataframe) Columns() []string {
	return c.headers
}

func (c *Dataframe) Rows() []map[string]*string {
	return c.rows
}

func CreateDataframe(records [][]string) (*Dataframe, error) {
	// if len = 1, just headers and no data
	if len(records) <= 1 {
		return nil, errors.New("no data provided for csv")
	}

	df := &Dataframe{
		// indexToCol: map[int]string{},
		// colToIndex: make(map[string]int, len(records[0])),
		// data:       make([][]string, len(records[0])),
		headers: make([]string, len(records[0])),
		rows:    make([]map[string]*string, len(records)-1),
	}

	df.headers = records[0]
	// indexToCol := map[int]string{}
	// colToEntries := map[string][]string{}

	set := map[string]bool{}
	// make sure no duplicate headers
	for _, header := range records[0] {
		_, ok := set[header]
		if ok {
			return nil, errors.New("duplicate header")
		}
		set[header] = true

		// csv.colToIndex[header] = idx
		// // preserve order of csv
		// indexToCol[idx] = header
		// colToEntries[header] = make([]string, len(records)-1)
		// csv.headers[idx] = header
		// // fmt.Println("SETTING HEADER ", idx, " TO: ", header)
	}
	// csv.colToEntries = colToEntries
	// now process the data

	for idx, row := range records[1:] {
		df.rows[idx] = map[string]*string{}

		for j, data := range row {
			someData := data
			header := df.headers[j]
			if data == "" {
				df.rows[idx][header] = nil
			} else {
				df.rows[idx][header] = &someData
			}
		}
		// fmt.Println("ROW")
		// for h, v := range csv.rows[idx] {
		// 	fmt.Println(h, *v)
		// }

		// fmt.Println(cols)

		// j is the index of the column header
		// for j, colValue := range cols {
		// 	csv.data[j] = append(csv.data[j], colValue)
		// }
		// for
	}

	// for _, h := range csv.headers {
	// 	fmt.Println(h)

	// 	entries := colToEntries[h]
	// 	fmt.Println(entries)
	// 	for _, v := range entries {
	// 		fmt.Println(*v)
	// 	}
	// }

	return df, nil
}
