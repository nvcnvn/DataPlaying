package main

import (
	"encoding/json"
	"github.com/nvcnvn/DataPlaying/data"
	"net/http"
)

func HandleAjaxClustering(c *context) {
	type requestObject struct {
		K      int
		Fields []int
	}
	var reqObj requestObject
	err := json.NewDecoder(c.Request.Body).Decode(&reqObj)
	n := len(reqObj.Fields)
	if err != nil || len(reqObj.Fields) == 0 {
		http.Error(c, "Invalid input array id", http.StatusBadRequest)
		return
	}

	fields := make([]data.DataField, n)
	for i := 0; i < n; i++ {
		fields[i] = DATASET.Data[reqObj.Fields[i]]
	}

	result, err := data.Clustering(reqObj.K, data.Real, fields...)
	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}
	if json.NewEncoder(c).Encode(&result) != nil {
		http.Error(c, "JSON marshal issue", http.StatusInternalServerError)
	}
}
