package main

import (
	"encoding/json"
	"net/http"
)

func HandleAjaxScatter(c *context) {
	var a []int
	err := json.NewDecoder(c.Request.Body).Decode(&a)
	if err != nil || len(a) != 2 {
		http.Error(c, "Invalid input array id", http.StatusBadRequest)
		return
	}

	n := len(DATASET.Data[a[0]].Real)
	result := make([][2]float64, 0, n)

	for i := 0; i < n; i++ {
		result = append(result, [2]float64{
			DATASET.Data[a[0]].Real[i],
			DATASET.Data[a[1]].Real[i],
		})
	}

	if json.NewEncoder(c).Encode(&result) != nil {
		http.Error(c, "JSON marshal issue", http.StatusInternalServerError)
	}
}
