package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LvlError = "ERROR"
	LvlWarn  = "WARN"
	LvlInfo  = "INFO"
	LvlDebug = "DEBUG"
	LvlTrace = "TRACE"
)

// New created a logger with given level
func New(level string) *Logger {
	return NewWithWriter(level, os.Stdout)
}

// NewWithWriter created a logger with given level and writer
func NewWithWriter(levelParam string, writer io.Writer) *Logger {
	level := MustGetValidLevel(levelParam)
	return &Logger{
		level:        level,
		writer:       writer,
		debugEnabled: level == LvlDebug || level == LvlTrace,
		traceEnabled: level == LvlTrace,
	}
}

type Logger struct {
	writer       io.Writer
	level        string
	debugEnabled bool
	traceEnabled bool
	mutex        sync.Mutex
}

// GetLevel returns the level in a thread safe way
func (l *Logger) GetLevel() string {
	return l.level
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.debugEnabled
}

// IsTraceEnabled returns true if trace logging is enabled
func (l *Logger) IsTraceEnabled() bool {
	return l.traceEnabled
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.logw(LvlError, msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logw(LvlWarn, msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logw(LvlInfo, msg, keysAndValues...)
}

// Debug should be used for detailed logs
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.debugEnabled {
		l.logw(LvlDebug, msg, keysAndValues...)
	}
}

// Trace should be used for dumps of payloads or similar
func (l *Logger) Trace(msg string, keysAndValues ...interface{}) {
	if l.traceEnabled {
		l.logw(LvlTrace, msg, keysAndValues...)
	}
}

// MustGetValidLevel returns a valid level or panics.
func MustGetValidLevel(level string) string {
	level, err := GetValidLevel(level)
	if err != nil {
		panic(err)
	}
	return level
}

// GetValidLevel parses a level string and returns a valid level name if found.
func GetValidLevel(level string) (string, error) {
	var allLevels = []string{LvlError, LvlWarn, LvlInfo, LvlDebug, LvlTrace}
	for _, l := range allLevels {
		if strings.EqualFold(l, level) {
			return l, nil
		}
	}
	return "", fmt.Errorf("invalid level: %s", level)
}

func (l *Logger) logw(level string, message string, keysAndValues ...interface{}) {
	// We must lock here, because we don't know for sure if the current io.writer uses locking
	l.mutex.Lock()
	defer l.mutex.Unlock()

	sw := MakeStackWriter(l.writer)
	defer sw.Flush()

	now := formatTimeCustom(time.Now())

	sw.Write("{\"level\": \"")
	sw.WriteEscaped(level)
	sw.Write("\", \"ts\": \"")
	sw.Write(string(now[:]))
	sw.Write("\", \"msg\": \"")
	sw.WriteEscaped(message)
	sw.Write("\"")

	fn := len(keysAndValues)
	for i := 0; i+1 < fn; i += 2 {
		sw.Write(", \"")
		switch key := keysAndValues[i].(type) {
		case string:
			sw.WriteEscaped(key)
		default:
			//writeUnknownValue(&sw, key)
			sw.WriteEscaped(fmt.Sprintf("%s", key))
		}
		sw.Write("\": \"")

		switch value := keysAndValues[i+1].(type) {
		case string:
			sw.WriteEscaped(value)
		case int:
			sw.WriteEscaped(strconv.Itoa(value))
		default:
			//sw.WriteEscaped(fmt.Sprintf("%s", value))
		}
		sw.Write("\"")
	}

	if includeCallerInfo(level) {
		funcName, fileName, line := retrieveCallInfo()
		sw.Write("\", \"caller\": \"")
		sw.WriteEscaped(funcName)
		sw.Write("() ")
		sw.WriteEscaped(fileName)
		sw.Write(":")
		sw.Write(strconv.Itoa(line))
		sw.Write("\"")
	}

	sw.Write("}\n")
}

func writeUnknownValue(sw *StackWriter, value interface{}) {
	s := fmt.Sprintf("%d", 10)
	sw.WriteEscaped(s)
}

func includeCallerInfo(level string) bool {
	return level == LvlError || level == LvlWarn
}

func retrieveCallInfo() (funcName string, file string, line int) {
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return "", "", -1
	}
	return runtime.FuncForPC(pc).Name(), file, line
}
