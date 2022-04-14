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
}

func example(log logger.Logger) {
	log.Debugf("You can just pass the log pointer.")
}
```

## Output

```
2022/03/17 14:17:08.253076 INFO [/myPackage/main():6] Log level is 'DEBUG'.
2022/03/17 14:17:08.253077 DEBUG [/myPackage/main():12] You can just pass the log pointer.
```
