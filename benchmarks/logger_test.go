package benchmarks

import (
	"github.com/fond-of-vertigo/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"testing"
)

func BenchmarkLogger_Info(b *testing.B) {
	log := logger.NewWithWriter(logger.LvlInfo, io.Discard)
	longstring := makeString(50)
	for i := 0; i < b.N; i++ {
		log.Info("Lorem \"ipsum\"",
			"Key", longstring,
			"K2", 34875634,
			"K3", 1.25)
	}
}

func BenchmarkLogger_zap_Infow(b *testing.B) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(io.Discard), zap.InfoLevel)
	logger := zap.New(core).Sugar()
	defer logger.Sync()
	longstring := makeString(50)
	for i := 0; i < b.N; i++ {
		logger.Infow("Lorem \"ipsum\"",
			"Key", longstring,
			"K2", 34875634,
			"K3", 1.25)
	}
}

func makeString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('0' + (i % 10)))
	}
	return sb.String()
}
