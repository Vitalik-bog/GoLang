package logger

import (
    "log"
    "os"
)

type Logger struct {
    infoLog  *log.Logger
    errorLog *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
        errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func (l *Logger) Info(format string, args ...interface{}) {
    if len(args) > 0 {
        l.infoLog.Printf(format, args...)
    } else {
        l.infoLog.Println(format)
    }
}

func (l *Logger) Error(format string, args ...interface{}) {
    if len(args) > 0 {
        l.errorLog.Printf(format, args...)
    } else {
        l.errorLog.Println(format)
    }
}
