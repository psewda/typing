package log

import (
	"io"

	"github.com/labstack/gommon/log"
	"github.com/psewda/typing/internal/utils"
)

type (
	// LevelType is the log level type.
	LevelType uint8
)

const (
	// LevelTypeDebug is representing 'Debug' value for log level type.
	LevelTypeDebug LevelType = iota

	// LevelTypeInfo is representing 'Info' value for log level type.
	LevelTypeInfo

	// LevelTypeWarn is representing 'Warn' value for log level type.
	LevelTypeWarn

	// LevelTypeError is representing 'Error' value for log level type.
	LevelTypeError
)

const (
	// Header is the print format for log. It is used
	// to include/exclude the fields in json output.
	Header = `{"time":"${time_rfc3339}","level":"${level}"}`
)

// Configuration is the config info for initializing new logger.
type Configuration struct {
	// Level is the log level for logging data.
	Level LevelType

	// Output is the log writer such as Stdout, Buffer, File etc.
	Output io.Writer

	// Color is to enable/disable output in colored mode.
	Color bool
}

// Logger is the logger struct.
type Logger struct {
	logger *log.Logger
}

// Debug writes the logs on the configured output. If the log
// level is higher than DEBUG, log writing is ignored.
func (l *Logger) Debug(value string) {
	l.logger.Debug(value)
}

// Info writes the logs on the configured output. If the log
// level is higher than INFO, log writing is ignored.
func (l *Logger) Info(value string) {
	l.logger.Info(value)
}

// Warn writes the logs on the configured output. If the log
// level is higher than WARN, log writing is ignored.
func (l *Logger) Warn(value string) {
	l.logger.Warn(value)
}

// Error writes the logs on the configured output. If the log
// level is higher than ERROR, log writing is ignored.
func (l *Logger) Error(value string, err error) {
	l.logger.Error(utils.AppendError(value, err))
}

// Fatal writes the logs on the configured output and
// exits the process with 0 error code.
func (l *Logger) Fatal(value string, err error) {
	l.logger.Fatal(utils.AppendError(value, err))
}

// New creates a new instance of logger.
func New(config Configuration) *Logger {
	logr := log.New("-")
	logr.SetHeader(Header)

	switch config.Level {
	case LevelTypeDebug:
		logr.SetLevel(log.DEBUG)
	case LevelTypeInfo:
		logr.SetLevel(log.INFO)
	case LevelTypeWarn:
		logr.SetLevel(log.WARN)
	case LevelTypeError:
		logr.SetLevel(log.ERROR)
	}

	if config.Output != nil {
		logr.SetOutput(config.Output)
	}

	logr.EnableColor()
	if !config.Color {
		logr.DisableColor()
	}

	return &Logger{
		logger: logr,
	}
}
