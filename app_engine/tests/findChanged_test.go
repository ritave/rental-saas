package tests

import (
	"reflect"
	"testing"
	"calendar-synch/logic"
)

func TestCompare(t *testing.T) {
	type args struct {
		saved  logic.SortableEvents
		actual logic.SortableEvents
	}
	tests := []struct {
		name    string
		args    args
		want    []logic.EventModified
		wantErr bool
	}{
		{"empty", args{logic.SortableEvents{}, logic.SortableEvents{}}, nil, false},
		{},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logic.Compare(tt.args.saved, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
