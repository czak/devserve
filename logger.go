package main

import (
	"fmt"
)

var logger *simpleLogger = &simpleLogger{}

const (
	LevelDebug = -4
	LevelInfo  = 0
	LevelWarn  = 4
	LevelError = 8
)

type simpleLogger struct {
	level int
}

func (l *simpleLogger) Debug(format string, v ...any) {
	if l.level <= LevelDebug {
		fmt.Println("\033[90m[DEBUG]\033[0m", fmt.Sprintf(format, v...))
	}
}

func (l *simpleLogger) Info(format string, v ...any) {
	if l.level <= LevelInfo {
		fmt.Println("\033[34m[INFO]\033[0m ", fmt.Sprintf(format, v...))
	}
}

func (l *simpleLogger) Warn(format string, v ...any) {
	if l.level <= LevelWarn {
		fmt.Println("\033[33m[WARN]\033[0m ", fmt.Sprintf(format, v...))
	}
}

func (l *simpleLogger) Error(format string, v ...any) {
	if l.level <= LevelError {
		fmt.Println("\033[91m[ERROR]\033[0m", fmt.Sprintf(format, v...))
	}
}
