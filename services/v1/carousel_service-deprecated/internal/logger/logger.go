package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

// Automatically called at package import
func init() {
	Info = log.New(os.Stdout, "[INF] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WRN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "[DBG] ", log.Ldate|log.Ltime|log.Lshortfile)
}

type TraceData struct {
	File string
	Line int
	Fn   string
}

func (t *TraceData) ToString() string {
	return fmt.Sprintf("%s:%d:", t.Fn, t.Line)
}

func Trace() *TraceData {
	// func trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return &TraceData{"?", 0, "?"}
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return &TraceData{file, line, "?"}
	}
	return &TraceData{file, line, fn.Name()}
}
