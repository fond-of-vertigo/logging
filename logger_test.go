package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLogger_Info(t *testing.T) {
	logger := New(LvlInfo)
	for i := 0; i < 1000; i++ {
		logger.Info("Test %d: %s", i, "Lorem ipsum")
	}
}

/*func TestLogger_InfoW(t *testing.T) {
	logger := New(LvlInfo)
	for i := 0; i < 1000; i++ {
		logger.Infow("Lorem ipsum",
			"Test", i,
		)
	}

	if actualMsg.Message != "Test msg" {
		t.Errorf("Message is incorrect, Expected %s, Actual %s", "Test msg", actualMsg.Message)
	}

	if actualMsg.Value == "" {
		t.Errorf("Key `Key1` does not exist")
	}

	if actualMsg.Value != "Value1" {
		t.Errorf("Value is incorrect, Expected %s, Actual %s", "Value1", actualMsg.Value)
	}
}

func TestLogger_MultiThread(t *testing.T) {
	tests := []struct {
		count       int
		threadCount int
	}{
		{
			count:       100,
			threadCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.count), func(t *testing.T) {

			var wg sync.WaitGroup
			wg.Add(tt.threadCount) // threadCount routines we need to wait for.

			logger := New(LvlInfo)
			for i := 0; i < tt.threadCount; i++ {
				go func(i int) {
					defer wg.Done()
					for j := 0; j < tt.count; j++ {
						logger.Info("Test %d: %s", i*tt.threadCount+j, "Lorem ipsum")
					}
				}(i)
			}
			wg.Wait()
		})
	}
}

func TestLogger_Perf_Structured(t *testing.T) {
	tests := []struct {
		count int
	}{
		{count: 1000 * 1000},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.count), func(t *testing.T) {
			logger := NewWithWriter(LvlInfo, os.NewFile(0, os.DevNull))
			start := time.Now()
			for i := 0; i < tt.count; i++ {
				logger.Info("Lorem ipsum", "Test ", i)
			}
			elapsed := time.Since(start)
			fmt.Printf("Logging %d took %s\n", tt.count, elapsed)
		})
	}
}

func TestLogger_Allocs_Structured(t *testing.T) {
	logger := NewWithWriter(LvlInfo, os.Stdout)
	key := makeString(10)
	longstring := makeString(4096)
	allocs := testing.AllocsPerRun(1, func() {
		logger.Info("Lorem \"ipsum\"",
			key, longstring,
			"int", 1,
			"bool", true)
	})

	if allocs > 0.0 {
		t.Errorf("Allocs detected! Want 0 allocs, got %f", allocs)
	}
}
