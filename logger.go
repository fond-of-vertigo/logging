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

type Logger interface {
	Error(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Trace(msg string, keysAndValues ...interface{})

	GetLevel() string
	IsDebugEnabled() bool
	IsTraceEnabled() bool
	Named(n string) Logger
	clone() Logger
	setName(n string)
	getName() string
	// Levels INFO, WARN and ERROR are always enabled.
}

// New created a logger with given level
func New(level string) Logger {
	return NewWithWriter(level, os.Stdout)
}

// NewWithWriter created a logger with given level and writer
func NewWithWriter(levelParam string, writer io.Writer) Logger {
	level := MustGetValidLevel(levelParam)
	return &instance{
		level:        level,
		writer:       writer,
		debugEnabled: level == LvlDebug || level == LvlTrace,
		traceEnabled: level == LvlTrace,
	}
}

type instance struct {
	writer       io.Writer
	level        string
	debugEnabled bool
	traceEnabled bool
	mutex        sync.Mutex
	name         string
}

func (l *instance) setName(n string) {
	l.name = n
}

func (l *instance) getName() string {
	return l.name
}

func (l *instance) Named(n string) Logger {
	if n == "" {
		return l
	}
	log := l.clone()
	if log.getName() == "" {
		log.setName(n)
	} else {
		log.setName(strings.Join([]string{l.name, n}, "."))
	}
	return log
}

func (l *instance) clone() Logger {
	return &instance{
		level:        l.level,
		name:         l.name,
		writer:       l.writer,
		debugEnabled: l.IsDebugEnabled(),
		traceEnabled: l.IsTraceEnabled(),
	}
}

// GetLevel returns the level in a thread safe way
func (l *instance) GetLevel() string {
	return l.level
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *instance) IsDebugEnabled() bool {
	return l.debugEnabled
}

// IsTraceEnabled returns true if trace logging is enabled
func (l *instance) IsTraceEnabled() bool {
	return l.traceEnabled
}

func (l *instance) Error(msg string, keysAndValues ...interface{}) {
	l.log(LvlError, msg, keysAndValues...)
}

func (l *instance) Warn(msg string, keysAndValues ...interface{}) {
	l.log(LvlWarn, msg, keysAndValues...)
}

func (l *instance) Info(msg string, keysAndValues ...interface{}) {
	l.log(LvlInfo, msg, keysAndValues...)
}

// Debug should be used for detailed logs
func (l *instance) Debug(msg string, keysAndValues ...interface{}) {
	if l.debugEnabled {
		l.log(LvlDebug, msg, keysAndValues...)
	}
}

// Trace should be used for dumps of payloads or similar
func (l *instance) Trace(msg string, keysAndValues ...interface{}) {
	if l.traceEnabled {
		l.log(LvlTrace, msg, keysAndValues...)
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

func (l *instance) log(level string, message string, keysAndValues ...interface{}) {
	// We must lock here, because we don't know for sure if the current io.writer uses locking
	l.mutex.Lock()
	defer l.mutex.Unlock()

	sw := MakeStackWriter(l.writer)
	defer sw.Flush()

	now := FormatLogTime(time.Now())

	sw.Write("{\"ts\": ")
	sw.WriteJSONString(string(now[:]))
	sw.Write(", \"level\": ")
	sw.WriteJSONString(level)
	if l.name != "" {
		sw.Write(", \"logger\": ")
		sw.WriteJSONString(l.name)
	}
	sw.Write(", \"message\": ")
	sw.WriteJSONString(message)

	fn := len(keysAndValues)
	for i := 0; i+1 < fn; i += 2 {
		sw.Write(", ")
		encodeKey(&sw, noescape_interface(&keysAndValues[i]))
		sw.Write(": ")
		encodeValue(&sw, noescape_interface(&keysAndValues[i+1]))
	}

	if includeCallerInfo(level) {
		funcName, fileName, line := retrieveCallInfo()
		sw.Write(", \"caller_func\": ")
		sw.WriteJSONString(funcName)
		sw.Write(", \"caller_file\": \"")
		sw.WriteEscaped(fileName)
		sw.Write(":")
		sw.Write(strconv.Itoa(line))
		sw.Write("\"")
	}

	sw.Write("}\n")
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
