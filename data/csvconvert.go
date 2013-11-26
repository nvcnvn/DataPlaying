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

	numfields := len(records)
	header := make([]Attr, 0, numfields)
	data := make([]DataField, 0, numfields)
	sortedData := make([]DataField, 0, numfields)

	for i := 0; i < numfields; i++ {
		header = append(header, Attr{strconv.Itoa(i), Real})
		record := records[i]

		field := make(sort.Float64Slice, 0, len(record))
		for _, s := range record {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				field = append(field, f)
			} else {
				field = append(field, 0.0)
			}
		}

		sortField := make(sort.Float64Slice, 0, len(record))
		copy(field, sortField)
		sort.Sort(sortField)

		data = append(data, DataField{Real: field})
		sortedData = append(sortedData, DataField{Real: sortField})
	}

	return &DataSet{header, data, sortedData}, nil
}
