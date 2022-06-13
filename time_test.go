package logger

import (
	"testing"
	"time"
)

func TestFormatLogTime(t *testing.T) {
	tests := []struct {
		name       string
		timeString string
	}{{
		name:       "Equal time string as std lib",
		timeString: "2022-02-01T13:01:02.123456Z",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := time.Parse(time.RFC3339Nano, tt.timeString)
			if err != nil {
				t.Fatalf("Inavlid test data, invalid time string \"%s\": %s", tt.timeString, err.Error())
			}

			formattedTime := FormatLogTime(ts)
			ourTimeString := string(formattedTime[:])
			if ourTimeString != tt.timeString {
				t.Errorf("FormatLogTime() = %v, want %v", ourTimeString, tt.timeString)
			}

			stdLibTimeString := ts.UTC().Format(time.RFC3339Nano)
			if ourTimeString != stdLibTimeString {
				t.Errorf("FormatLogTime() = %v, want %v", ourTimeString, stdLibTimeString)
			}
		})
	}
}

func TestFormatLogTime_ZeroAlloc(t *testing.T) {
	now := time.Now()
	allocs := testing.AllocsPerRun(1, func() {
		ts := FormatLogTime(now)
		str := string(ts[:])
		if len(str) > 0 {
		}
	})

	if allocs > 0.0 {
		t.Errorf("Allocs detected! Want 0 allocs, got %f", allocs)
	}
}
