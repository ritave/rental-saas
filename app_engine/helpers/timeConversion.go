package helpers

import "time"

const defaultTimeType = time.RFC3339

//StringToTime assumes the input string is correct (very wild assumption, ik)
func StringToTime(in string) (time.Time) {
	out, _ := time.Parse(defaultTimeType, in)
	return out
}