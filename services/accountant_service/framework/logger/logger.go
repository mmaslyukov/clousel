package logger

import (
	"log"
	"os"

	"accountant_service/framework/core"
)

type Logger struct {
	e *log.Logger
	w *log.Logger
	i *log.Logger
	d *log.Logger
}

func (l *Logger) Err() *log.Logger {
	return l.e
}
func (l *Logger) Wrn() *log.Logger {
	return l.w
}
func (l *Logger) Inf() *log.Logger {
	return l.i
}
func (l *Logger) Dbg() *log.Logger {
	return l.d
}

func LoggerCreate() core.ILogger {
	return &Logger{
		e: log.New(os.Stdout, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile),
		w: log.New(os.Stdout, "[WRN] ", log.Ldate|log.Ltime|log.Lshortfile),
		i: log.New(os.Stdout, "[INF] ", log.Ldate|log.Ltime|log.Lshortfile),
		d: log.New(os.Stdout, "[DBG] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
