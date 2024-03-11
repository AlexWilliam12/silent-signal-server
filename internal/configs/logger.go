package configs

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

// Logger struct
type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	writer  io.Writer
}

// Get a instance of a new Logger
func NewLogger(prefix string) *Logger {

	// Logger instance
	writer := io.Writer(os.Stdout)
	logger := log.New(writer, prefix, log.Ldate|log.Ltime)

	// Colours for Logger options
	debug := color.New(color.FgCyan).SprintFunc()
	info := color.New(color.FgGreen).SprintFunc()
	warning := color.New(color.FgYellow).SprintFunc()
	err := color.New(color.FgRed).SprintFunc()

	return &Logger{
		debug:   log.New(writer, debug("DEBUG: "), logger.Flags()),
		info:    log.New(writer, info("INFO: "), logger.Flags()),
		warning: log.New(writer, warning("WARNING: "), logger.Flags()),
		err:     log.New(writer, err("ERROR: "), logger.Flags()),
		writer:  writer,
	}
}

// Print Debugger line
func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

// Print Information line
func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

// Print Warning line
func (l *Logger) Warning(v ...interface{}) {
	l.warning.Println(v...)
}

// Print Error line
func (l *Logger) Err(v ...interface{}) {
	l.err.Println(v...)
}

// Print Debugger line with formatted string
func (l *Logger) Debugf(str string, v ...interface{}) {
	l.debug.Printf(str, v...)
}

// Print Information line with formatted string
func (l *Logger) Infof(str string, v ...interface{}) {
	l.info.Printf(str, v...)
}

// Print Warning line with formatted string
func (l *Logger) Warningf(str string, v ...interface{}) {
	l.warning.Printf(str, v...)
}

// Print Error line with formatted string
func (l *Logger) Errf(str string, v ...interface{}) {
	l.err.Printf(str, v...)
}
