package main

import (
	"testing"
	"time"

	"snippetbox.micypac.io/internal/assert"
)

// func TestHumanDate(t *testing.T) {
// 	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
// 	got := humanDate(tm)
// 	want := "17 Mar 2022 at 10:15"

// 	if got != want {
// 		t.Errorf("got %q, want %q", got, want)
// 	}
// }

func TestHumanDate(t *testing.T) {
	tests := []struct{
		name string
		tm time.Time
		want string 
	}{
		{
			name: "UTC",
			tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.tm)

			assert.Equal(t, got, tt.want)
		})
	}
}
