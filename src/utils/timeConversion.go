package utils

import (
	"time"
)

const DefaultTimeType = time.RFC3339
const thousand = 1000
const million = thousand*thousand

func VerifyStringToTime(in string) (time.Time, error) {
	return time.Parse(DefaultTimeType, in)
}

//StringToTime assumes the input string is correct (very wild assumption, ik)
func StringToTime(in string) (time.Time) {
	out, _ := time.Parse(DefaultTimeType, in)
	return out
}

func TimeToString(in time.Time) (string) {
	return in.Format(time.RFC3339)
}

func MillisecondsToTime(in int64) (time.Time) {
	milliRemainder := in%thousand
	nanoseconds := million * milliRemainder

	return time.Unix(in/thousand, nanoseconds)
}

func TimeToMilliseconds(in time.Time) (int64) {
	return in.UnixNano()/million
}