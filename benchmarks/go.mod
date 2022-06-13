module github.com/fond-of-vertigo/logger/benchmarks

go 1.18

require (
	github.com/fond-of-vertigo/logger v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.27.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace github.com/fond-of-vertigo/logger => ../
