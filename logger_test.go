package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLogger_Infof(t *testing.T) {
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
}*/

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
				logger.Info("Test %d: %s", i, "Lorem ipsum")
			}
			elapsed := time.Since(start)
			fmt.Printf("Logging %d took %s\n", tt.count, elapsed)
		})
	}
}

func TestLogger_Perf_Structured(t *testing.T) {
	/*tests := []struct {
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
	}*/
}

func TestLogger_Allocs_Structured(t *testing.T) {
	logger := NewWithWriter(LvlInfo, os.Stdout)
	key := makeString(10)
	longstring := makeString(32768)
	allocs := testing.AllocsPerRun(1, func() {
		logger.Info("Lorem \"ipsum\"", key, longstring)
	})
	fmt.Printf("\nAllocs: %f\n", allocs)
}

func BenchmarkLogger_Info(b *testing.B) {
	logger := NewWithWriter(LvlInfo, io.Discard)
	longstring := makeString(50)
	for i := 0; i < b.N; i++ {
		logger.Info("Lorem \"ipsum\"",
			"Key", longstring,
			"K2", 34875634,
			"K3", "sdfjiosdfjio")
	}
}

func TestLogger_zap_Allocs_Structured(t *testing.T) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(os.Stdout), zap.InfoLevel)
	logger := zap.New(core)
	longstring := makeString(1024)
	allocs := testing.AllocsPerRun(1, func() {
		logger.Info("Lorem \"ipsum\"",
			zap.String("Key", longstring),
		)
	})
	logger.Sync()
	fmt.Printf("\nAllocs: %f\n", allocs)
}

func BenchmarkLogger_zap_Infow(b *testing.B) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(io.Discard), zap.InfoLevel)
	logger := zap.New(core)
	defer logger.Sync()
	longstring := makeString(50)
	for i := 0; i < b.N; i++ {
		logger.Info("Lorem \"ipsum\"",
			zap.String("Key", longstring),
		)
	}
}

func BenchmarkLogger_time_std(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		now.Format("2006-01-02 15:04:05.999999")
	}
}

func BenchmarkLogger_time_cust(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		formatTimeCustom(now)
	}
}

func BenchmarkLogger_time_now(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}
