package data

import (
	"encoding/csv"
	"io"
	"sort"
	"strconv"
)

func CSVConvert(r io.Reader) (*DataSet, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}

	var header []Attr

	titleRow := records[0]
	numfields := len(titleRow)
	numRows := len(records)
	header = make([]Attr, 0, numfields)
	data := make([]DataField, 0, numfields)
	sortedData := make([]DataField, 0, numfields)

	for _, v := range titleRow {
		header = append(header, Attr{Name: v, Type: Real})
		data = append(data, DataField{
			Real: make(sort.Float64Slice, numRows-1),
		})
		sortedData = append(sortedData, DataField{
			Real: make(sort.Float64Slice, numRows-1),
		})
	}

	for i := 1; i < numRows; i++ {
		record := records[i]

		var (
			idx int
			s   string
		)
		for idx, s = range record {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				data[idx].Real[i-1] = f
				sortedData[idx].Real[i-1] = f
			}
		}
	}

	for i := 0; i < numfields; i++ {
		sort.Sort(sortedData[i].Real)
	}

	return &DataSet{header, data, sortedData}, nil
}
