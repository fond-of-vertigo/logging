# Logging library

Logging will be outputted to the stdout.

## Example

```go
package myPackage

import "github.com/fond-of-vertigo/logging"

func main() {
	log := logger.New(logger.LvlDebug)
	log.Infof("Log level is '%s'.", log.GetLevel())

	example(log)
	
	// Structured logging
	log.Infow("Log message",
		"key1", "value1",
		"key2", "value2")
}

func example(log logger.Logger) {
	log.Debugf("You can just pass the log pointer.")
}
```

## Output

```
2022/03/17 14:17:08.253076 INFO [/myPackage/main():7] Log level is 'DEBUG'.
2022/03/17 14:17:08.253077 DEBUG [/myPackage/example():18] You can just pass the log pointer.
{"level": "INFO", "time": "2022/03/17 14:17:08.253080", "caller": "/myPackage/main():12", "message": "Log message", 
"key1": "value1", "key2": "value2"}
```
