package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
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
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Tracef(format string, v ...interface{})

	Infow(message string, keysAndValues ...interface{})

	GetLevel() string
	IsDebugEnabled() bool
	IsTraceEnabled() bool
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
		logger:       log.New(writer, "", log.Ldate|log.Ltime|log.Lmicroseconds),
		debugEnabled: level == LvlDebug || level == LvlTrace,
		traceEnabled: level == LvlTrace,
	}
}

type instance struct {
	logger       *log.Logger
	level        string
	debugEnabled bool
	traceEnabled bool
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

// Errorf for errors, prints to stderr
func (l *instance) Errorf(format string, v ...interface{}) {
	l.logf(LvlError, format, v...)
}

// Warnf for default log messages
func (l *instance) Warnf(format string, v ...interface{}) {
	l.logf(LvlWarn, format, v...)
}

// Infof for default log messages
func (l *instance) Infof(format string, v ...interface{}) {
	l.logf(LvlInfo, format, v...)
}

// Infow for structured log messages
func (l *instance) Infow(message string, keysAndValues ...interface{}) {
	l.logw(LvlInfo, message, keysAndValues...)
}

// Debugf should be used for detailed logs
func (l *instance) Debugf(format string, v ...interface{}) {
	if l.debugEnabled {
		l.logf(LvlDebug, format, v...)
	}
}

// Tracef should be used for dumps of payloads or similar
func (l *instance) Tracef(format string, v ...interface{}) {
	if l.traceEnabled {
		l.logf(LvlTrace, format, v...)
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

// logf for default log messages of given level
func (l *instance) logf(level string, format string, v ...interface{}) {
	l.logger.Printf(l.prependMetadata(level, format), v...)
}

// logw for structured log messages of given level
func (l *instance) logw(level string, message string, keysAndValues ...interface{}) {
	now := time.Now()
	nowFormatted := now.Format("2006-01-02 15:04:05.999999")
	pkgName, funcName, line := retrieveCallInfo()
	callerInfo := fmt.Sprintf("%s/%s():%d", pkgName, funcName, line)
	if err := l.write("{" +
		"\"level\": \"" + level + "\"" +
		", \"time\": \"" + nowFormatted + "\"" +
		", \"caller\": \"" + callerInfo + "\"" +
		", \"message\": \"" + message + "\""); err != nil {
		return
	}
	for i := 0; i < len(keysAndValues); i++ {
		key, ok := keysAndValues[i].(string)
		i++
		if !ok {
			continue
		}
		value := fmt.Sprintf("%v", keysAndValues[i])

		if err := l.write(", \"" + key + "\": \"" + value + "\""); err != nil {
			return
		}
	}
	if err := l.write("}\n"); err != nil {
		return
	}
}

func (l *instance) write(text string) error {
	_, err := l.logger.Writer().Write([]byte(text))
	return err
}

func (l *instance) prependMetadata(levelName string, str string) string {
	pkgName, funcName, line := retrieveCallInfo()
	return fmt.Sprintf("%s [%s/%s():%d] %s", levelName, pkgName, funcName, line, str)
}

func retrieveCallInfo() (pkgName string, funcName string, line int) {
	pc, _, line, ok := runtime.Caller(4)
	if !ok {
		return "", "", -1
	}
	funcPath := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcPath, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcPath[lastSlash:], '.') + lastSlash

	return funcPath[lastSlash:lastDot], funcPath[lastDot+1:], line
}
