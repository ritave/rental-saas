package utils

import (
	"time"
)

const DefaultTimeType = time.RFC3339
const thousand = 1000
const million = thousand*thousand
var CET = time.FixedZone("CET", 3600)

func VerifyStringToTime(in string) (time.Time, error) {
	return time.Parse(DefaultTimeType, in)
}

//StringToTime assumes the input string is correct (very wild assumption, ik)
func StringToTime(in string) (time.Time) {
	out, _ := time.Parse(DefaultTimeType, in)
	return out
}

func TimeToString(in time.Time) (string) {
	return in.In(CET).Format(time.RFC3339)
}

func MillisecondsToTime(in int64) (time.Time) {
	milliRemainder := in%thousand
	nanoseconds := million * milliRemainder

	return time.Unix(in/thousand, nanoseconds)
}

func TimeToMilliseconds(in time.Time) (int64) {
	return in.UnixNano()/million
}

const POZAMIATANE_DateFormat = "2006-01-02"
//'documentation' suggest they are using this format: '2018-06-01 12:00:00'
func POZAMIATANE_DatetimeToDateString(in time.Time) (string) {
	return in.In(CET).Format(POZAMIATANE_DateFormat)
}

const POZAMIATANE_TimeFormat = "15:04:05"
//'documentation' suggest they are using this format: '2018-06-01 12:00:00'
func POZAMIATANE_DatetimeToTimeString(in time.Time) (string) {
	return in.In(CET).Format(POZAMIATANE_TimeFormat)
}
