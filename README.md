# Logging library

Logging will be outputted to the stdout.

## Basic Example

```go
package myPackage

import "github.com/fond-of-vertigo/logger"

func main() {
	log := logger.New(logger.LvlDebug)

	example(log)

	// Structured logging
	log.Info("Log message",
		"key1", "value1",
		"key2", "value2")
}

func example(log logger.Logger) {
	log.Debug("You can just pass the log pointer.")
}

// out
// {"time": "2022/03/17 14:17:08.253080", "level": "INFO", "message": "You can just pass the log pointer."}
// {"time": "2022/03/17 14:17:08.253080", "level": "INFO", "message": "Log message", "key1": "value1", "key2": "value2"}
```

## Named Logger

Logger can also have names.

```go
package myPackage

import "github.com/fond-of-vertigo/logger"

func main() {
	log := logger.New(logger.LvlDebug).Named("main")
	log_sub := log.Named("sub")

	log.Info("Log message", "key1", "value1")
	log_sub.Info("Log message2", "key2", "value2")
}

// out
// {"time": "2022/03/17 14:17:08.253080", "level": "INFO", "logger": "main", "message": "Log message", "key1": "value1", "key2": "value2"}
// {"time": "2022/03/17 14:17:08.253080", "level": "INFO", "logger": "main.sub", "message": "Log message2", "key2": "value2", "key2": "value2"}

```



