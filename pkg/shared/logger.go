package shared

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags|log.Lshortfile),
		Error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		Debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Infof(format string, args ...any)  { l.Info.Printf(format, args...) }
func (l *Logger) Errorf(format string, args ...any) { l.Error.Printf(format, args...) }
func (l *Logger) Debugf(format string, args ...any) { l.Debug.Printf(format, args...) }
