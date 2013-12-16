package data

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
)

func LibSVMConvert(r io.Reader) (*DataSet, error) {
	scanner := bufio.NewScanner(r)

	type pair struct {
		key int
		val float64
	}

	set := make([][]pair, 0, 10)
	size := 0
	for scanner.Scan() {
		/* scanner.Bytes() should return a line.
		words will be a slice ["-1", "a:b", "c:d"]
		*/
		words := strings.Split(scanner.Text(), " ")
		line := make([]pair, 0, len(words)-1)
		for _, w := range words[1:] {
			part := strings.Split(w, ":")

			key, err := strconv.Atoi(part[0])
			if err != nil {
				continue
			}

			val, _ := strconv.ParseFloat(part[1], 64)
			line = append(line, pair{key, val})
		}

		if len(line) == 0 {
			continue
		}

		if l := line[len(line)-1].key; l > size {
			size = l
		}

		set = append(set, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	numrows := len(set)
	var dset DataSet
	dset.Header = make([]Attr, 0, size)
	dset.Data = make([]DataField, 0, size)
	dset.SortedData = make([]DataField, 0, size)

	for i := 0; i < size; i++ {
		dset.Header = append(dset.Header, Attr{
			"Attr " + strconv.Itoa(i+1),
			Real,
		})

		dset.Data = append(dset.Data, DataField{
			Real: make(sort.Float64Slice, numrows, numrows),
		})
		dset.SortedData = append(dset.SortedData, DataField{
			Real: make(sort.Float64Slice, numrows, numrows),
		})
	}

	for i, line := range set {
		for _, p := range line {
			dset.Data[p.key-1].Real[i] = p.val
			dset.SortedData[p.key-1].Real[i] = p.val

		}
	}

	for i := 0; i < size; i++ {
		sort.Sort(dset.SortedData[i].Real)
	}
	return &dset, nil
}
