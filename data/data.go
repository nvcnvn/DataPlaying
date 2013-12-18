package data

import (
	"errors"
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

// Clustering with K-means
func Clustering(k int, t AttrType, a ...DataField) ([][]int, error) {
	// just implement for Real data
	rows := len(a[0].Real)
	if k > rows {
		return nil, errors.New("k bigger than number of records")
	}

	nAttrs := len(a)
	//select k centroids, not smart way...just fast
	centroids := make([][]float64, k) // centroids vector
	r := rows / k
	for i := 0; i < k; i++ {
		v := make([]float64, 0, nAttrs)
		for _, f := range a {
			v = append(v, f.Real[i*r])
		}
		centroids[i] = v
	}

	// run until the centroids has no change
	for {
		mCluster := make([][]int, k) // cluster lookup map
		for i := 0; i < k; i++ {
			mCluster[i] = make([]int, 0, r)
		}

		record := make([]float64, nAttrs)
		// choose a cluster for a record
		for i := 0; i < rows; i++ {
			for j := 0; j < nAttrs; j++ {
				record[j] = a[j].Real[i]
			}

			// get the first distance
			var currentK int
			var dist, tmpDist float64
			var tmp float64
			for j := 0; j < nAttrs; j++ {
				x := record[j] - centroids[0][j]
				tmp += x * x
			}

			dist = math.Sqrt(tmp)

			for t := 1; t < k; t++ {
				tmp = 0.0
				for j := 0; j < nAttrs; j++ {
					x := record[j] - centroids[t][j]
					tmp += x * x
				}
				tmpDist = math.Sqrt(tmp)
				if tmpDist < dist {
					currentK = t
					dist = tmpDist
				}
			}
			// now record i in cluster currentK
			mCluster[currentK] = append(mCluster[currentK], i)
		}
		// re-calculate centroids
		tmpCentroids := make([][]float64, k)
		for i := 0; i < k; i++ {
			tmpCentroids[i] = make([]float64, nAttrs)
			for j := 0; j < nAttrs; j++ {
				var x float64
				for _, v := range mCluster[i] {
					x += a[j].Real[v]
				}
				if len(mCluster[i]) > 0 {
					tmpCentroids[i][j] = x / float64(len(mCluster[i]))
				}
			}
		}
		// compare just calculated and the old centroids
		var hasDiff bool
		for i := 0; i < k; i++ {
			for j := 0; j < nAttrs; j++ {
				if centroids[i][j] != tmpCentroids[i][j] {
					hasDiff = true
					break
				}
			}
		}

		if hasDiff {
			centroids = tmpCentroids
		} else {
			return mCluster, nil
		}
	}

	return nil, errors.New("Program error")
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

// GetMedian must recieve a sorted DataField
func GetMedian(a DataField, t AttrType) float64 {
	var median float64

	switch t {
	case Real:
		n := len(a.Real)
		if n%2 == 0 {
			n = n / 2
			median = (a.Real[n-1] + a.Real[n]) / 2.0
		} else {
			n = n / 2
			median = a.Real[n]
		}
	case Integer:
		n := len(a.Integer)
		if n%2 == 0 {
			n = n / 2
			median = float64(a.Integer[n-1]+a.Integer[n]) / 2.0
		} else {
			n = n / 2
			median = float64(a.Integer[n])
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
		for _, v := range a.Integer {
			lookup[v]++
		}
	case String:
		for _, v := range a.String {
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
		n := len(a.Real)
		if n == 1 {
			v = 0.0
			break
		}
		var sum float64
		for i := 0; i < n; i++ {
			v := a.Real[i] - mean
			sum += v * v
		}
		v = sum / float64(n-1)
	case Integer:
		n := len(a.Integer)
		if n == 1 {
			v = 0.0
			break
		}
		var sum float64
		for i := 0; i < n; i++ {
			v := float64(a.Integer[i]) - mean
			sum += v * v
		}
		v = sum / float64(n-1)
	}

	return v
}

// GetQuartiles return a slice of num-Quantiles, the input
// a DataField must be sorted
func GetQuantiles(a DataField, t AttrType, num int) []float64 {
	quartiles := make([]float64, 0, num+1)
	switch t {
	case Real:
		n := len(a.Real)
		quartiles = append(quartiles, a.Real[0])
		for i := 1; i < num; i++ {
			realIndex := float64(i) / float64(num) * float64(n-1)
			index := math.Floor(realIndex)
			frac := realIndex - index
			if index+1.0 < float64(n) {
				quartiles = append(quartiles, a.Real[int(index)]*
					(1-frac)+a.Real[int(index)+1]*frac)
			} else {
				quartiles = append(quartiles, a.Real[int(index)])
			}
		}
		quartiles = append(quartiles, a.Real[len(a.Real)-1])
	case Integer:
		n := len(a.Integer)
		quartiles = append(quartiles, float64(a.Integer[0]))
		for i := 1; i < num; i++ {
			realIndex := float64(i) / float64(num) * float64(n-1)
			index := math.Floor(realIndex)
			frac := realIndex - index
			if index+1.0 < float64(n) {
				quartiles = append(quartiles, float64(a.Integer[int(index)])*
					(1-frac)+float64(a.Integer[int(index)+1])*frac)
			} else {
				quartiles = append(quartiles, float64(a.Integer[int(index)]))
			}
		}
		quartiles = append(quartiles, float64(a.Integer[len(a.Integer)-1]))
	}

	return quartiles
}
