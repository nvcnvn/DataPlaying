package data

import (
	"math"
	"sort"
)

type AttrType int

const (
	Invalid AttrType = iota
	Numeric
	Real
	Integer
	String
	Nominal
	Ordinal
)

type Attr struct {
	Name string
	Type AttrType
}

type DataField struct {
	Real    sort.Float64Slice
	Integer sort.IntSlice
	String  []string
	Generic []interface{}
}

type DataSet struct {
	Header     []Attr
	Data       []DataField
	SortedData []DataField
}

func GetMean(a DataField, t AttrType) float64 {
	var mean float64

	switch t {
	case Real:
		n := float64(len(a.Real))
		var sum float64
		for _, v := range a.Real {
			sum += v
		}
		mean = sum / n
	case Integer:
		n := float64(len(a.Integer))
		var sum int
		for _, v := range a.Integer {
			sum += v
		}
		mean = float64(sum) / n
	}

	return mean
}

func GetMedian(a DataField, t AttrType) float64 {
	var median float64

	switch t {
	case Real:
		n := len(a.Real)
		if n%2 == 0 {
			n = n / 2
			median = (a.Real[n] + a.Real[n+1]) / 2.0
		} else {
			n = n / 2
			median = a.Real[n+1] / 2.0
		}
	case Integer:
		n := len(a.Integer)
		if n%2 == 0 {
			n = n / 2
			median = (a.Real[n] + a.Real[n+1]) / 2.0
		} else {
			n = n / 2
			median = a.Real[n+1] / 2.0
		}
	}

	return median
}

// A data structure to hold a key/value pair.
type pair struct {
	key   interface{}
	value int
}

// A slice of pairs that implements sort.Interface to sort by value.
type pairList []pair

func (p pairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p pairList) Len() int {
	return len(p)
}

func (p pairList) Less(i, j int) bool {
	return p[i].value < p[j].value
}

func GetPresentCount(a DataField, t AttrType) map[interface{}]int {
	lookup := make(map[interface{}]int)
	switch t {
	case Real:
		for _, v := range a.Real {
			lookup[v]++
		}
	case Integer:
		for _, v := range a.Real {
			lookup[v]++
		}
	case String:
		for _, v := range a.Real {
			lookup[v]++
		}
	default:
		for _, v := range a.Generic {
			lookup[v]++
		}
	}
	return lookup
}

func GetMode(a DataField, t AttrType, lookup map[interface{}]int) []interface{} {
	lst := make(pairList, 0, len(lookup))
	for key, val := range lookup {
		lst = append(lst, pair{key, val})
	}
	sort.Sort(sort.Reverse(lst))

	result := make([]interface{}, 0, 1)

	for _, v := range lst {
		if v.value != lst[0].value {
			break
		}
		result = append(result, v.key)
	}

	return result
}

func GetVariance(a DataField, t AttrType, mean float64) float64 {
	var v float64

	switch t {
	case Real:
		var sum float64
		for _, v := range a.Real {
			sum += math.Pow(v, 2)
		}
		v = sum*0.5 - math.Pow(mean, 2)
	case Integer:
		var sum int
		for _, v := range a.Integer {
			sum += v * v
		}
		v = float64(sum)*0.5 - math.Pow(mean, 2)
	}

	return v
}
