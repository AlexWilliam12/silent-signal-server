package configs

// import (
// 	"io"
// 	"log"
// 	"os"

// 	"github.com/fatih/color"
// )

// type Logger struct {
// 	debug   *log.Logger
// 	info    *log.Logger
// 	warning *log.Logger
// 	err     *log.Logger
// 	writer  io.Writer
// }

// func NewLogger(prefix string) *Logger {

// 	writer := io.Writer(os.Stdout)
// 	logger := log.New(writer, prefix, log.Ldate|log.Ltime)

// 	debug := color.New(color.FgCyan).SprintFunc()
// 	info := color.New(color.FgGreen).SprintFunc()
// 	warning := color.New(color.FgYellow).SprintFunc()
// 	err := color.New(color.FgRed).SprintFunc()

// 	return &Logger{
// 		debug:   log.New(writer, debug("DEBUG: "), logger.Flags()),
// 		info:    log.New(writer, info("INFO: "), logger.Flags()),
// 		warning: log.New(writer, warning("WARNING: "), logger.Flags()),
// 		err:     log.New(writer, err("ERROR: "), logger.Flags()),
// 		writer:  writer,
// 	}
// }

// func (l *Logger) Debug(v ...interface{}) {
// 	l.debug.Println(v...)
// }

// func (l *Logger) Info(v ...interface{}) {
// 	l.info.Println(v...)
// }

// func (l *Logger) Warning(v ...interface{}) {
// 	l.warning.Println(v...)
// }

// func (l *Logger) Err(v ...interface{}) {
// 	l.err.Println(v...)
// }

// func (l *Logger) Debugf(str string, v ...interface{}) {
// 	l.debug.Printf(str, v...)
// }

// func (l *Logger) Infof(str string, v ...interface{}) {
// 	l.info.Printf(str, v...)
// }

// func (l *Logger) Warningf(str string, v ...interface{}) {
// 	l.warning.Printf(str, v...)
// }

// func (l *Logger) Errf(str string, v ...interface{}) {
// 	l.err.Printf(str, v...)
// }
