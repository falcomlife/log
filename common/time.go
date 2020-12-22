package common

import (
	"strconv"
	"time"
)

func GetRangeTime(period int64) (string, string) {
	timeUnix := time.Now().Unix()
	timeperiod := timeUnix - period
	return strconv.FormatInt(timeUnix, 10), strconv.FormatInt(timeperiod, 10)
}
