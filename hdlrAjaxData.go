package main

import (
	"encoding/json"
	"net/http"
)

func HandleAjaxData(c *context) {
	var a []int
	err := json.NewDecoder(c.Request.Body).Decode(&a)
	n := len(a)
	if err != nil || len(a) == 0 {
		http.Error(c, "Invalid input array id", http.StatusBadRequest)
		return
	}

	type NamedData struct {
		Name string
		Data []float64
	}

	result := make([]NamedData, 0, n)

	for i := 0; i < n; i++ {
		result = append(result, NamedData{
			DATASET.Header[a[i]].Name,
			DATASET.Data[a[i]].Real,
		})
	}

	if json.NewEncoder(c).Encode(&result) != nil {
		http.Error(c, "JSON marshal issue", http.StatusInternalServerError)
	}
}

func HandleAjaxSortedData(c *context) {
	var a []int
	err := json.NewDecoder(c.Request.Body).Decode(&a)
	n := len(a)
	if err != nil || len(a) == 0 {
		http.Error(c, "Invalid input array id", http.StatusBadRequest)
		return
	}

	type NamedData struct {
		Name string
		Data []float64
	}

	result := make([]NamedData, 0, n)

	for i := 0; i < n; i++ {
		result = append(result, NamedData{
			DATASET.Header[a[i]].Name,
			DATASET.SortedData[a[i]].Real,
		})
	}

	if json.NewEncoder(c).Encode(&result) != nil {
		http.Error(c, "JSON marshal issue", http.StatusInternalServerError)
	}
}
