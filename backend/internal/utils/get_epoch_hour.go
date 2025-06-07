package utils

import "time"

func GetEpochHour(t time.Time) int64 {
	return t.Unix() / 3600
}

func GetTimeFromEpochHour(hr int64) time.Time {
	return time.Unix(hr*3600, 0)
}
