package common

import (
	"fmt"
	"k8s.io/klog/v2"
	"sort"
	"strconv"
)

// get top one percent value in a float arr
func Median(arr []float64) float64 {
	sort.Sort(sort.Reverse(sort.Float64Slice(arr)))
	count := len(arr)
	return (arr[((count)/100)+1] + arr[((count)/100)]) / 2
}

// ellipsis dot
func EllipsisDot(f float64) float64 {
	float, err := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	if err != nil {
		klog.Error(err)
		return 0
	}
	return float
}
