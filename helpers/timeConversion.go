package helpers

import "time"

const DefaultTimeType = time.RFC3339

//StringToTime assumes the input string is correct (very wild assumption, ik)
func StringToTime(in string) (time.Time) {
	out, _ := time.Parse(DefaultTimeType, in)
	return out
}

func TimeToString(in time.Time) (string) {
	return in.Format(time.RFC3339)
}