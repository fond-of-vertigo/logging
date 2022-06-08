package benchmarks

import (
	"github.com/fond-of-vertigo/logger"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"testing"
)

func BenchmarkLogger_Info(b *testing.B) {
	log := logger.NewWithWriter(logger.LvlInfo, io.Discard)
	longstring := makeString(50)
	alloc := testing.AllocsPerRun(b.N, func() {
		log.Info("Lorem \"ipsum\"",
			"Key", longstring,
			"K2", 34875634,
			"K3", 1.25)
	})
	b.Logf("Allocations:  %f", alloc)

}

func BenchmarkLogger_zap_Infow(b *testing.B) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(io.Discard), zap.InfoLevel)
	logger := zap.New(core).Sugar()
	defer logger.Sync()
	longstring := makeString(50)
	alloc := testing.AllocsPerRun(b.N, func() {
		logger.Infow("Lorem \"ipsum\"",
			"Key", longstring,
			"K2", 34875634,
			"K3", 1.25)
	})
	b.Logf("Allocations:  %f", alloc)
}

func BenchmarkLogger_zerolog_Info(b *testing.B) {
	log := zerolog.New(io.Discard).With().Timestamp().Logger()
	longstring := makeString(50)
	alloc := testing.AllocsPerRun(b.N, func() {
		log.Info().Str("Key", longstring).Int("K2", 34875634).Float64("K3", 1.25).Msg("Lorem \"ipsum\"")
	})
	b.Logf("Allocations:  %f", alloc)
}

func makeString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('0' + (i % 10)))
	}
	return sb.String()
}
