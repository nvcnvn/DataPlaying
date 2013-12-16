package main

import (
	"encoding/json"
	"github.com/nvcnvn/DataPlaying/data"
	"log"
	"net/http"
	"path"
)

var DATASET *data.DataSet

func HandleIndex(c *context) {
	viewdata := c.ViewData("title")

	if c.Post("submit", false) == "Send" {
		file, fHeader, err := c.Request.FormFile("file")
		if err != nil {
			http.Error(c, "Internal Error",
				http.StatusInternalServerError)
			return
		}

		var set *data.DataSet

		if path.Ext(fHeader.Filename) == ".libsvm" {
			set, err = data.LibSVMConvert(file)
		} else {
			set, err = data.CSVConvert(file)
		}

		if err == nil {
			viewdata["dataset"] = set
			DATASET = set
		} else {
			http.Error(c, "Invalid data format",
				http.StatusBadRequest)
			log.Println("HandleIndex", err)
		}
	}

	c.View("index_detail.tmpl", viewdata)
}

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

func HandleAjax(c *context) {
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
