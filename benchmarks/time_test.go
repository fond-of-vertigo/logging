package benchmarks

import (
	"github.com/fond-of-vertigo/logger"
	"testing"
	"time"
)

func BenchmarkLogger_FormatLogTime(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		logger.FormatLogTime(now)
	}
}

func BenchmarkLogger_FormatTimeStdLib(b *testing.B) {
	now := time.Now().UTC()
	for i := 0; i < b.N; i++ {
		now.String()
	}
}
