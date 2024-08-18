package config

import (
	"log"
	"os"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	file  *os.File
}

var logPath = "./.log"

func NewLogger() *Logger {
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	return &Logger{
		file:  f,
		debug: log.New(f, " [DEBUG] ", log.LstdFlags),
		info:  log.New(f, " [INFO]  ", log.LstdFlags),
		warn:  log.New(f, " [WARN]  ", log.LstdFlags),
		err:   log.New(f, " [ERROR] ", log.LstdFlags),
	}
}

// Close fecha o arquivo de log.
func (l *Logger) CloseLogger(msg *string) error {
	os.Truncate(logPath, 0)
	return l.file.Close()
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.warn.Println(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.warn.Printf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.err.Println(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.err.Printf(format, v...)
}
