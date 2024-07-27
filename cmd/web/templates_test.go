package main

import (
	"kibonga/quickbits/internal/assert"
	"testing"
	"time"
)

const T1 string = "27 Jul 2024 at 09:59"

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 7, 27, 10, 5, 0, 0, time.UTC),
			want: "27 Jul 2024 at 10:05",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 7, 27, 10, 5, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "27 Jul 2024 at 09:05",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			hd := humanDate(ts.tm)
			assert.Equal(t, hd, ts.want)
		})
	}
}
