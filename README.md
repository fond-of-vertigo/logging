# Logging library

Logging will be outputted to the stdout.

## Example

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
// {"level": "INFO", "time": "2022/03/17 14:17:08.253080", "message": "You can just pass the log pointer."}
// {"level": "INFO", "time": "2022/03/17 14:17:08.253080", "message": "Log message", "key1": "value1", "key2": "value2"}
```





