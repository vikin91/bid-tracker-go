package logging

import (
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var once sync.Once

var instance *log.Logger

//Logger returns a global singleton logger object to access go-kit logger
func Logger() *log.Logger {
	once.Do(func() {
		l := ConfigureLogger(level.AllowDebug())
		instance = &l
	})
	return instance
}

//ConfigureLogger sets the log level
func ConfigureLogger(logLevel level.Option) log.Logger {
	w := log.NewSyncWriter(os.Stdout)
	logger := log.NewLogfmtLogger(w)
	logger = level.NewFilter(logger, logLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.Caller(4))
	return logger
}

// LogError logs an error with the singleton logger with message and error
func LogError(message string, err error) {
	level.Error(*Logger()).Log("message", message, "error", err)
}

// LogWarning logs a warning with the singleton logger with message and error
func LogWarning(message string, err error) {
	level.Warn(*Logger()).Log("message", message, "error", err)
}

// LogInfo logs an info with the singleton logger with message and error
func LogInfo(message string) {
	level.Info(*Logger()).Log("message", message)
}
