package logger

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLogger_Infof(t *testing.T) {
	logger := New(LvlInfo)
	for i := 0; i < 1000; i++ {
		logger.Infof("Test %d: %s", i, "Lorem ipsum")
	}
}

func TestLogger_InfoW(t *testing.T) {
	logger := New(LvlInfo)
	for i := 0; i < 1000; i++ {
		logger.Infow("Lorem ipsum",
			"Test", i,
		)
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
						logger.Infof("Test %d: %s", i*tt.threadCount+j, "Lorem ipsum")
					}
				}(i)
			}
			wg.Wait()
		})
	}
}

func TestLogger_Perf(t *testing.T) {
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
				logger.Infof("Test %d: %s", i, "Lorem ipsum")
			}
			elapsed := time.Since(start)
			fmt.Printf("Logging %d took %s\n", tt.count, elapsed)
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
				logger.Infow("Lorem ipsum", "Test ", i)
			}
			elapsed := time.Since(start)
			fmt.Printf("Logging %d took %s\n", tt.count, elapsed)
		})
	}
}
