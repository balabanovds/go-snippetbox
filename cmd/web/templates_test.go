package main

import (
	"testing"
	"time"
)

func TestHumanTime(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 07, 16, 3, 0, 0, 0, time.UTC),
			want: "16 Jul 2020 at 03:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "MSK",
			tm:   time.Date(1982, 07, 16, 6, 0, 0, 0, time.FixedZone("MSK", 3*60*60)),
			want: "16 Jul 1982 at 03:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.tm)

			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}
