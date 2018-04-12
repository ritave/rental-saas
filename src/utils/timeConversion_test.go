package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestVerifyStringToTime(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"UTC", args{"2017-07-13T13:07:00Z"}, time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC), false},
		{"CET", args{"2017-07-13T13:07:00+01:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CET), false},
		{"CEST", args{"2017-07-13T13:07:00+02:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CEST), false},
		// Fun fact: this test fails when given as want: time.Date(..., local) where local, _ := time.LoadLocation("Poland")
		{"Invalid", args{"2017-07-13T13:07:00Z+02:00"}, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyStringToTime(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyStringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("VerifyStringToTime() = %v, want %v", got, tt.want)
				t.Errorf("(Unix) Got: %d, want: %d", got.Unix(), tt.want.Unix())
			}
		})
	}
}

func TestStringToTime(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"UTC", args{"2017-07-13T13:07:00Z"}, time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)},
		{"CET", args{"2017-07-13T13:07:00+01:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CET)},
		{"CEST", args{"2017-07-13T13:07:00+02:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CEST)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToTime(tt.args.in); !got.Equal(tt.want) {
				t.Errorf("(%s) StringToTime() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTimeToString(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "2017-07-13T14:07:00+01:00"},
		{"CET", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CET)}, "2017-07-13T13:07:00+01:00"},
	}

	// this is stupid
	timeTravel.To(CET, CETNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToString(tt.args.in); got != tt.want {
				t.Errorf("TimeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeToString2(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "2017-07-13T15:07:00+02:00"},
		{"CEST", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CEST)}, "2017-07-13T13:07:00+02:00"},
	}

	// this is stupid
	timeTravel.To(CEST, CESTNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToString(tt.args.in); got != tt.want {
				t.Errorf("TimeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMillisecondsToTime(t *testing.T) {
	type args struct {
		in int64
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"zero", args{0}, time.Unix(0, 0)},
		{"one millisecond", args{1}, time.Unix(0, million)},
		{"almost a second", args{999}, time.Unix(0, 999*million)},
		{"second in nano", args{1000}, time.Unix(0, thousand*million)},
		{"just a second", args{1000}, time.Unix(1, 0)},
		{"second + millisecond", args{1001}, time.Unix(1, million)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MillisecondsToTime(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(%s) MillisecondsToTime() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTimeToMilliseconds(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"zero", args{time.Unix(0, 0)}, 0},
		{"1 nanosecond", args{time.Unix(0, 1)}, 0},
		{"almost a millisecond of nanoseconds", args{time.Unix(0, 999*thousand)}, 0},
		{"millisecond of nanoseconds", args{time.Unix(0, thousand*thousand)}, 1},
		{"second in nanoseconds", args{time.Unix(0, thousand*million)}, 1000},
		{"just a second", args{time.Unix(1, 0)}, 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToMilliseconds(tt.args.in); got != tt.want {
				t.Errorf("(%s) TimeToMilliseconds() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_DatetimeToDateString(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "2017-07-13"},
		{"CET", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CET)}, "2017-07-13"},
		{"UTC", args{time.Date(2017, 7, 13, 23, 59, 0, 0, time.UTC)}, "2017-07-14"}, // ^^ cheeky
		{"CET", args{time.Date(2017, 7, 13, 23, 59, 0, 0, CET)}, "2017-07-13"},
	}

	timeTravel.To(CET, CETNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POZAMIATANE_DatetimeToDateString(tt.args.in); got != tt.want {
				t.Errorf("DateTimeToDateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_DatetimeToTimeString(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "14:07:00"},
		{"CET", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CET)}, "13:07:00"},
		{"UTC", args{time.Date(2017, 7, 13, 23, 59, 0, 0, time.UTC)}, "00:59:00"}, // ^^ cheeky
		{"CET", args{time.Date(2017, 7, 13, 23, 59, 0, 0, CET)}, "23:59:00"},
	}

	timeTravel.To(CET, CETNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POZAMIATANE_DatetimeToTimeString(tt.args.in); got != tt.want {
				t.Errorf("DateTimeToTimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_DatetimeToTimeStringCEST(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "15:07:00"},
		{"CEST", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CEST)}, "13:07:00"},
		{"UTC", args{time.Date(2017, 7, 13, 23, 59, 0, 0, time.UTC)}, "01:59:00"}, // ^^ cheeky
		{"CEST", args{time.Date(2017, 7, 13, 23, 59, 0, 0, CEST)}, "23:59:00"},
	}

	timeTravel.To(CEST, CESTNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POZAMIATANE_DatetimeToTimeString(tt.args.in); got != tt.want {
				t.Errorf("DateTimeToTimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_DatetimeToDatetimeString(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "2017-07-13 14:07:00"},
		{"CET", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CET)}, "2017-07-13 13:07:00"},
		{"UTC", args{time.Date(2017, 7, 13, 23, 59, 0, 0, time.UTC)}, "2017-07-14 00:59:00"}, // ^^ cheeky
		{"CET", args{time.Date(2017, 7, 13, 23, 59, 0, 0, CET)}, "2017-07-13 23:59:00"},
	}
	
	timeTravel.To(CET, CETNumeric)
	defer timeTravel.Back()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POZAMIATANE_DatetimeToDatetimeString(tt.args.in); got != tt.want {
				t.Errorf("DateTimeToDatetimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_DatetimeToDatetimeStringCEST(t *testing.T) {
	type args struct {
		in time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"UTC", args{time.Date(2017, 7, 13, 13, 7, 0, 0, time.UTC)}, "2017-07-13 15:07:00"},
		{"CEST", args{time.Date(2017, 7, 13, 13, 7, 0, 0, CEST)}, "2017-07-13 13:07:00"},
		{"UTC", args{time.Date(2017, 7, 13, 23, 59, 0, 0, time.UTC)}, "2017-07-14 01:59:00"}, // ^^ cheeky
		{"CEST", args{time.Date(2017, 7, 13, 23, 59, 0, 0, CEST)}, "2017-07-13 23:59:00"},
	}

	timeTravel.To(CEST, CESTNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := POZAMIATANE_DatetimeToDatetimeString(tt.args.in); got != tt.want {
				t.Errorf("DateTimeToDatetimeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_StringToDatetime(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"CET local", args{"2017-07-13 13:07:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CET), false},
		{"invalid format", args{"2017-07-13T13:07:00"}, NullDate, true},
		{"invalid date", args{"2017-99-99 13:07:00"}, NullDate, true},
		{"timezone supplied", args{"2017-99-99 13:07:00-0700"}, NullDate, true},
		{"null", args{""}, NullDate, true},
	}

	timeTravel.To(CET, CETNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := POZAMIATANE_StringToDatetimeLOCAL(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("POZAMIATANE_StringToDatetimeLOCAL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("POZAMIATANE_StringToDatetimeLOCAL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPOZAMIATANE_StringToDatetimeCEST(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"CEST local", args{"2017-07-13 13:07:00"}, time.Date(2017, 7, 13, 13, 7, 0, 0, CEST), false},
		{"invalid format", args{"2017-07-13T13:07:00"}, NullDate, true},
		{"invalid date", args{"2017-99-99 13:07:00"}, NullDate, true},
		{"timezone supplied", args{"2017-99-99 13:07:00-0700"}, NullDate, true},
		{"null", args{""}, NullDate, true},
	}

	timeTravel.To(CEST, CESTNumeric)
	defer timeTravel.Back()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := POZAMIATANE_StringToDatetimeLOCAL(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("POZAMIATANE_StringToDatetimeLOCAL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("POZAMIATANE_StringToDatetimeLOCAL() = %v, want %v", got, tt.want)
			}
		})
	}
}
