package common

import "sort"

// get top one percent value in a float arr
func Median(arr []float64) float64 {
	sort.Sort(sort.Reverse(sort.Float64Slice(arr)))
	count := len(arr)
	return (arr[((count)/100)-1] + arr[((count)/100)]) / 2
}
