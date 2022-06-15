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

type logMsg struct {
	Timestamp string `json:"ts"`
	Level     string `json:"level"`
	Name      string `json:"logger"`
	Message   string `json:"message"`
	Value     string `json:"Key1"`
}

func TestLogger_Info(t *testing.T) {
	logger := New(LvlInfo)
	for i := 0; i < 1000; i++ {
		logger.Info("Test %d: %s", i, "Lorem ipsum")
	}
}

func TestLogger_Info_CheckOutput(t *testing.T) {
	parentTestLoggerOut := bytes.NewBufferString("")
	parentTestLogger := NewWithWriter(LvlInfo, parentTestLoggerOut).Named("main")
	tests := []struct {
		testName               string
		logLevel               string
		loggerName             string
		expectedFullLoggerName string
		message                string
		key                    string
		value                  string
		parentLogger           Logger
		loggerBuffer           *bytes.Buffer
	}{
		{
			testName: "Simple structured log",
			logLevel: LvlInfo,
			message:  "test message",
			key:      "key1",
			value:    "value1",
		}, {
			testName:               "Simple named logger",
			logLevel:               LvlInfo,
			loggerName:             "main",
			expectedFullLoggerName: "main",
			message:                "test message",
			key:                    "key1",
			value:                  "value2",
		}, {
			testName:               "Nested named logger",
			logLevel:               LvlInfo,
			loggerName:             "sub",
			expectedFullLoggerName: "main.sub",
			message:                "test message",
			key:                    "key1",
			value:                  "value2",
			parentLogger:           parentTestLogger,
			loggerBuffer:           parentTestLoggerOut,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			logger := tt.parentLogger
			out := tt.loggerBuffer
			if logger == nil {
				out = bytes.NewBufferString("")
				logger = NewWithWriter(tt.logLevel, out)
			}
			logger = logger.Named(tt.loggerName)
			switch tt.logLevel {
			case LvlInfo:
				logger.Info(tt.message, tt.key, tt.value)
			case LvlDebug:
				logger.Debug(tt.message, tt.key, tt.value)
			case LvlError:
				logger.Error(tt.message, tt.key, tt.value)
			case LvlWarn:
				logger.Warn(tt.message, tt.key, tt.value)
			case LvlTrace:
				logger.Trace(tt.message, tt.key, tt.value)
			}

			actualMsg := logMsg{}
			if err := json.Unmarshal(out.Bytes(), &actualMsg); err != nil {
				t.Error(err)
			}

			if len(actualMsg.Timestamp) <= 0 {
				t.Errorf("Timestamp cannot be empty.")
			}

			if actualMsg.Level != tt.logLevel {
				t.Errorf("Level is incorrect, Expected %s, Actual %s", tt.logLevel, actualMsg.Level)
			}

			if logger.GetLevel() != tt.logLevel {
				t.Errorf("Level is incorrect, Expected %s, Actual %s", tt.logLevel, logger.GetLevel())
			}

			if actualMsg.Message != tt.message {
				t.Errorf("Message is incorrect, Expected %s, Actual %s", tt.message, actualMsg.Message)
			}

			if actualMsg.Value == "" {
				t.Errorf("Key `%s` does not exist", tt.key)
			}

			if actualMsg.Value != tt.value {
				t.Errorf("Value is incorrect, Expected %s, Actual %s", tt.value, actualMsg.Value)
			}

			if actualMsg.Name != tt.expectedFullLoggerName {
				t.Errorf("Logger should have name `%s`, but name is `%s`", tt.expectedFullLoggerName, actualMsg.Name)
			}
		})
	}
}

func TestABC(t *testing.T) {
	logger := NewWithWriter(LvlError, os.Stdout).Named("main")

	logger_t := logger.Named("sub")
	logger.Info("Test")
	logger.Info("Test", "k", "v")
	logger_t.Info("Test2", "kk", "vv")
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
