package main

import (
	"encoding/json"
	"github.com/nvcnvn/DataPlaying/data"
	"net/http"
)

type Lookup struct {
	Value []interface{}
	Count []int
}

type SumaryResponse struct {
	Min, Max, Mean, Median, Variance float64
	Quartiles                        []float64
	Mode                             []interface{}
	Name                             string
	Lookup                           Lookup
	Data                             []float64
}

func HandleAjaxSummary(c *context) {
	var a []int
	err := json.NewDecoder(c.Request.Body).Decode(&a)
	if err != nil {
		http.Error(c, "Invalid input array id", http.StatusBadRequest)
		return
	}

	result := make([]SumaryResponse, 0, len(a))
	for _, v := range a {
		if 0 <= v && v < len(DATASET.Header) {
			set := DATASET.Data[v]
			sortedSet := DATASET.SortedData[v]
			var (
				quartiles = data.GetQuantiles(sortedSet, data.Real, 4)
				lookupMap = data.GetPresentCount(sortedSet, data.Real)
				mean      = data.GetMean(set, data.Real)
				lookup    Lookup
			)

			lookup.Count = make([]int, 0, len(lookupMap))
			lookup.Value = make([]interface{}, 0, len(lookupMap))

			for k, v := range lookupMap {
				lookup.Value = append(lookup.Value, k)
				lookup.Count = append(lookup.Count, v)
			}

			result = append(result, SumaryResponse{
				Min:       sortedSet.Real[0],
				Max:       sortedSet.Real[len(sortedSet.Real)-1],
				Mean:      mean,
				Name:      DATASET.Header[v].Name,
				Median:    quartiles[2],
				Mode:      data.GetMode(set, data.Real, lookupMap),
				Variance:  data.GetVariance(set, data.Real, mean),
				Quartiles: quartiles,
				Lookup:    lookup,
				Data:      sortedSet.Real,
			})
		}
	}
	if json.NewEncoder(c).Encode(&result) != nil {
		http.Error(c, "JSON marshal issue", http.StatusInternalServerError)
	}
}
