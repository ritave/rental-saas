package utils

import (
	"time"
	"sync"
)

const DefaultTimeType = time.RFC3339
const thousand = 1000
const million = thousand*thousand

var CET = time.FixedZone("currentTZ", 3600)
const CETNumeric = "+0100"

var CEST = time.FixedZone("currentTZ", 7200)
const CESTNumeric = "+0200"

var NullDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

// FIXME this will lead to a bug on 28-10-2018
// but I'm leaving this as it is because, AFAIK, we may get rid of time changes in Poland
var currentTZ = CEST
var currentTZNumeric = CESTNumeric

var timeTravel = TimeTravel{}
type TimeTravel struct {
	before *time.Location
	mutex sync.Mutex
}
func (tt *TimeTravel) To(t *time.Location) {
	tt.mutex.Lock()
	tt.before = currentTZ
	currentTZ = t
}
func (tt *TimeTravel) Back() {
	currentTZ = tt.before
	tt.mutex.Unlock()
}

func VerifyStringToTime(in string) (time.Time, error) {
	return time.Parse(DefaultTimeType, in)
}

//StringToTime assumes the input string is correct (very wild assumption, ik)
func StringToTime(in string) (time.Time) {
	out, _ := time.Parse(DefaultTimeType, in)
	return out
}

func TimeToString(in time.Time) (string) {
	return in.In(currentTZ).Format(time.RFC3339)
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
	return in.In(currentTZ).Format(POZAMIATANE_DateFormat)
}

const POZAMIATANE_TimeFormat = "15:04:05"
//'documentation' suggest they are using this format: '2018-06-01 12:00:00'
func POZAMIATANE_DatetimeToTimeString(in time.Time) (string) {
	return in.In(currentTZ).Format(POZAMIATANE_TimeFormat)
}

const POZAMIATANE_DatetimeFormat = POZAMIATANE_DateFormat + " " + POZAMIATANE_TimeFormat
func POZAMIATANE_DatetimeToDatetimeString(in time.Time) (string) {
	return in.In(currentTZ).Format(POZAMIATANE_DatetimeFormat)
}

const POZAMIATANE_DatetimeFormatCET = POZAMIATANE_DatetimeFormat + "-0700"
func POZAMIATANE_StringToDatetimeLOCAL(in string) (time.Time, error) {
	//return time.Parse(POZAMIATANE_DatetimeFormat, in)
	return time.Parse(POZAMIATANE_DatetimeFormatCET, in +currentTZNumeric)
}
